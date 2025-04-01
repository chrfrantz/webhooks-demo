package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	structs "webhooks/webhooks"
)

// Initialize signature (via init())
var SignatureKey = "X-SIGNATURE"

// var Mac hash.Hash
var Secret []byte

// Webhook DB
var webhooks = []structs.WebhookRegistration{}

/*
Handles webhook registration (POST) and lookup (GET) requests.
Expects WebhookRegistration struct body in request.
*/
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Expects incoming body in terms of WebhookRegistration struct
		webhook := structs.WebhookRegistration{}
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusBadRequest)
		}
		webhooks = append(webhooks, webhook)
		// Note: Approach does not guarantee persistence or permanence of resource id (for CRUD)
		log.Println("Webhook " + webhook.Url + " has been registered.")
		// Print index of recorded webhook as response - note: in practice you would return some unique identifier, not exposing DB internals
		http.Error(w, strconv.Itoa(len(webhooks)-1), http.StatusCreated)
	case http.MethodGet:
		// For now just return all webhooks, don't respond to specific resource requests
		err := json.NewEncoder(w).Encode(webhooks)
		if err != nil {
			http.Error(w, "Something went wrong: "+err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "Method "+r.Method+" not supported for "+structs.WebhookEndpoint, http.StatusMethodNotAllowed)
	}
}

/*
Invokes the web service to trigger event. Currently only responds to POST requests.
*/
func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	str, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error during decoding message content. Error: " + string(str))
	}
	switch r.Method {
	case http.MethodPost:
		log.Println("Received POST request...")
		// Iterate through registered webhooks and invoke based on registered URL, method, and with received content
		for _, v := range webhooks {
			log.Println("Trigger event: Call to service endpoint for event " + v.Event +
				" and content '" + string(str) + "'.")
			go CallUrl(v.Url, v.Event, string(str))
		}
	default:
		http.Error(w, "Method "+r.Method+" not supported for "+structs.ServiceEndpoint, http.StatusMethodNotAllowed)
	}
}

/*
Calls given URL with given event indicator as well as content, and awaits response (status and body).
*/
func CallUrl(url string, event string, content string) {
	log.Println("Attempting invocation of url " + url + " with content '" + content + "'.")

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte("Event "+event+
		" occurred. Payload: "+content)))
	if err != nil {
		log.Printf("%v", "Error during request creation. Error:", err)
		return
	}

	/// BEGIN: HEADER GENERATION FOR CONTENT-BASED VALIDATION

	// Hash content (for content-based validation; not relevant for URL-based validation)
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(content))
	if err != nil {
		log.Printf("%v", "Error during content hashing. Error:", err)
		return
	}
	// Convert hash to string & add to header to transport to client for validation
	req.Header.Add(SignatureKey, hex.EncodeToString(mac.Sum(nil)))

	/// END: CONTENT-BASED VALIDATION

	// Perform invocation
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request. Error:", err)
		return
	}

	// Read the response
	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response. Error:", err)
		return
	}

	log.Println("Webhook " + url + " invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}

/*
Default handler for server side displaying service information.
*/
func DefaultServerHandler(w http.ResponseWriter, r *http.Request) {

	// Define content type, so browser renders links correctly
	w.Header().Add("content-type", "text/html")

	// Prepare output returned to client
	output := "This service offers the following endpoints: <br>The " + structs.WebhookEndpoint + " endpoint provides the registration functionality for webhooks (" + http.MethodPost + " and " + http.MethodGet + "), <br>The " +
		structs.ServiceEndpoint + " endpoint triggers the invocation of registered webhooks when called (with arbitrary payload)." +
		"<br>The payload structure for the webhook registration via " +
		http.MethodPost + " is the following JSON structure: {\"url\": \"http://targetHost:targetPort/pathTobeInvoked\", \"event\": \"EXAMPLE_EVENT\"}<br>" +
		"Please see the associated Readme for more information."

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}

/*
Default handler for client side displaying service information.
*/
func DefaultClientHandler(w http.ResponseWriter, r *http.Request) {

	// Define content type, so browser renders links correctly
	w.Header().Add("content-type", "text/html")

	// Prepare output returned to client
	output := "This service reacts on the following endpoint: <br>The /pathToBeInvoked endpoint reacts to invocation with any payload, " +
		"but can variably be configured to perform integrity checks (assessible via console output only)" +
		"<br>Please see the associated Readme for more information."

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
