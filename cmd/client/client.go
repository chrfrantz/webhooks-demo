package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"webhooks/handlers"
)

/*
Switch to toggle form of validation
Validation level 0: no validation; everything accepted
Validation level 1: check that URL is correct (obscured URL suffix)
Validation level 2: check that content is correctly encoded (does not check URL)
*/
var validationLevel = 0

// Invoked Hash to be accepted
var secret = []byte{1, 2, 3, 4, 5}        // not a good secret!
var urlMac = hmac.New(sha256.New, secret) // used for URL-based validation
var ClientSignatureKey = "X-SIGNATURE"    // used for content-based validation

// Sleep duration in seconds
var sleep = 0

// Endpoint on client
const ClientEndpoint = "/pathToBeInvoked"

/*
Dummy handler printing everything it receives to console and
confirm receipt to requester.
*/
func NonValidatingHandler(w http.ResponseWriter, r *http.Request) {

	// Sleep (if specified)
	if sleep > 0 {
		log.Println("Sleeping for " + strconv.Itoa(sleep) + " seconds ...")
	}
	time.Sleep(time.Duration(sleep) * time.Second)

	// Simply print body
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body: " + err.Error())
		http.Error(w, "Error when reading body: "+err.Error(), http.StatusBadRequest)
	}

	log.Println("Received invocation with method " + r.Method + " and body: " + string(content))

	// Writing response (Alternative: http.Error() function)
	_, err = fmt.Fprint(w, "Successfully invoked dummy web service.")
	if err != nil {
		log.Println("Something went wrong when sending response: " + err.Error())
	}
}

/*
Dummy handler printing everything it receives to console and checks
whether URL is correctly encoded.
*/
func URLValidatingHandler(w http.ResponseWriter, r *http.Request) {

	// Sleep (if specified)
	if sleep > 0 {
		log.Println("Sleeping for " + strconv.Itoa(sleep) + " seconds ...")
	}
	time.Sleep(time.Duration(sleep) * time.Second)

	// Simply print body
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body: " + err.Error())
		http.Error(w, "Error when reading body: "+err.Error(), http.StatusBadRequest)
	}

	log.Println("Received invocation with method " + r.Method + " and body: " + string(content))

	// Extract hash from URL
	split := strings.Split(r.URL.Path, "/")

	if len(split) != 3 {
		log.Println("Wrong number of tokens in " + r.URL.Path)
		http.Error(w, "Invalid invocation", http.StatusBadRequest)
		return
	}

	// Convert string to []byte
	received, err := hex.DecodeString(split[2])
	if err != nil {
		http.Error(w, "Error during hex decoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Invoked obscured URL component: ", split[2])

	// Compare HMAC with received request
	if hmac.Equal(received, urlMac.Sum(nil)) {
		log.Println("Valid invocation on " + r.URL.Path)
		_, err = fmt.Fprint(w, "Successfully invoked dummy web service.")
		if err != nil {
			log.Println("Something went wrong when sending response: " + err.Error())
		}
	} else { // Error - invalid HMAC
		log.Println("Call to non-existent webhook on " + r.URL.Path)
		http.Error(w, "Invalid invocation", http.StatusBadRequest)
	}
}

/*
Dummy handler printing everything it receives to console and checks
whether content is correctly encoded (with signature).
Note: The hash is reinitialized for each interaction.
Suggestion: Retain hash instance and write each invocation to it -
ensures integrity for all interactions
*/
func ContentValidatingHandler(w http.ResponseWriter, r *http.Request) {

	// Sleep (if specified)
	if sleep > 0 {
		log.Println("Sleeping for " + strconv.Itoa(sleep) + " seconds ...")
	}
	time.Sleep(time.Duration(sleep) * time.Second)

	// Simply print body
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error when reading body: " + err.Error())
		http.Error(w, "Error when reading body: "+err.Error(), http.StatusBadRequest)
	}

	log.Println("Received invocation with method " + r.Method + " and body: " + string(content))

	// Extract signature from header based on known key
	signature := r.Header.Get(ClientSignatureKey)

	// Convert string to []byte
	signatureByte, err := hex.DecodeString(signature)
	if err != nil {
		http.Error(w, "Error during Signature decoding: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Signature: " + signature)
	// Hash content of body
	mac := hmac.New(sha256.New, secret)
	_, err = mac.Write(content)
	if err != nil {
		http.Error(w, "Error during message decoding: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Content: " + hex.EncodeToString(mac.Sum(nil)))

	// Compare HMAC with received request
	if hmac.Equal(signatureByte, mac.Sum(nil)) {
		log.Println("Valid invocation (with validated content) on " + r.URL.Path)
		_, err = fmt.Fprint(w, "Successfully invoked dummy web service.")
		if err != nil {
			log.Println("Something went wrong when sending response: " + err.Error())
		}
	} else { // Error - invalid HMAC
		log.Println("Invalid invocation (tampered content?) on " + r.URL.Path)
		http.Error(w, "Invalid invocation", http.StatusBadRequest)
	}
}

func main() {

	port := "8081"

	// Environment variable constant for PaaS support
	PORT := "PORT"

	if os.Getenv(PORT) != "" {
		port = os.Getenv(PORT)
	}

	// Check for port in command line arguments (overrides defaults)
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) == 3 {
		// Read port number
		port = os.Args[1]
		// Read sleep parameter
		s, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("Error reading sleep parameter from command line (Value: " + os.Args[2] + ").")
		}
		sleep = s
	}

	log.Println("Service listening on port " + port)
	if sleep > 0 {
		log.Println("Activated execution delay on client (" + strconv.Itoa(sleep) + " seconds).")
	}
	// Register default handler
	http.HandleFunc("/", handlers.DefaultClientHandler)

	// Register functionality-specific handler
	switch validationLevel {
	case 0:
		// Register regular endpoint handler without any validation
		log.Println("Service URL (non-validating): http://localhost:" + port + ClientEndpoint)
		http.HandleFunc(ClientEndpoint, NonValidatingHandler)
	case 1:
		// Prepare URL-based validation based on obscured endpoint URL
		// Generate unique suffix (here based on secret, but other approaches conceivable, e.g., UUID)
		urlSuffix := hex.EncodeToString(urlMac.Sum(nil))
		log.Println("Service URL (URL-validating): http://localhost:" + port + ClientEndpoint + "/" + urlSuffix)
		// Append URL suffix during instantiation
		http.HandleFunc(ClientEndpoint+"/"+urlSuffix, URLValidatingHandler)
	case 2:
		// Register regular endpoint, but point to content-validating handler
		log.Println("Service URL (content-validating): http://localhost:" + port + ClientEndpoint)
		http.HandleFunc(ClientEndpoint, ContentValidatingHandler)
	default:
		log.Fatal("Invalid validation level. Exiting ...")
	}
	log.Println("Note: The URL above is the one that needs to be registered on the server's webhooks endpoint.")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
