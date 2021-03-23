package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
Switch to toggle on validation
Validation level 0: no validation; everything accepted
Validation level 1: check that URL is correct (signature)
Validation level 2: check that content is correctly encoded (does not check URL)
 */

var validationLevel = 0

// Invoked Hash to be accepted
var secret = []byte{1, 2, 3, 4, 5}        // not a good secret!
var urlMac = hmac.New(sha256.New, secret) // used for URL-based validation
var ClientSignatureKey = "X-SIGNATURE"    // used for content-based validation


/*
	Dummy handler printing everything it receives to console and
	confirm receipt to requester.
*/
func NonValidatingHandler(w http.ResponseWriter, r *http.Request) {

	// Simply print body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error when reading body: " + err.Error())
		http.Error(w, "Error when reading body: "+err.Error(), http.StatusBadRequest)
	}

	fmt.Println("Received invocation with method " + r.Method + " and body: " + string(content))

	// Writing response (Alternative: http.Error() function)
	_, err = fmt.Fprint(w, "Successfully invoked dummy web service.")
	if err != nil {
		fmt.Println("Something went wrong when sending response: " + err.Error())
	}
}

/*
	Dummy handler printing everything it receives to console and checks
	whether URL is correctly encoded.
 */
func URLValidatingHandler(w http.ResponseWriter, r *http.Request) {

	// Simply print body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error when reading body: " + err.Error())
		http.Error(w, "Error when reading body: " + err.Error(), http.StatusBadRequest)
	}

	fmt.Println("Received invocation with method " + r.Method + " and body: " + string(content))

	// Extract hash from URL
	split := strings.Split(r.URL.Path, "/")

	if len(split) != 3 {
		fmt.Println("Wrong number of tokens in " + r.URL.Path)
		http.Error(w, "Invalid invocation", http.StatusBadRequest)
		return
	}

	// Convert string to []byte
	received, err := hex.DecodeString(split[2])
	if err != nil {
		http.Error(w, "Error during HMAC decoding: " + err.Error(), http.StatusInternalServerError)
		return
	}

	// Compare HMAC with received request
	if hmac.Equal(received, urlMac.Sum(nil)) {
		fmt.Println("Valid invocation on " + r.URL.Path)
		_, err = fmt.Fprint(w, "Successfully invoked dummy web service.")
		if err != nil {
			fmt.Println("Something went wrong when sending response: " + err.Error())
		}
	} else { // Error - invalid HMAC
		fmt.Println("Call to non-existent webhook on " + r.URL.Path)
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

	// Simply print body
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error when reading body: " + err.Error())
		http.Error(w, "Error when reading body: " + err.Error(), http.StatusBadRequest)
	}

	fmt.Println("Received invocation with method " + r.Method + " and body: " + string(content))

	// Extract signature from header based on known key
	signature := r.Header.Get(ClientSignatureKey)

	// Convert string to []byte
	signatureByte, err := hex.DecodeString(signature)
	if err != nil {
		http.Error(w, "Error during Signature decoding: " + err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Signature: " + signature)
	// Hash content of body
	mac := hmac.New(sha256.New, secret)
	_, err = mac.Write(content)
	if err != nil {
		http.Error(w, "Error during message decoding: " + err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Content: " + hex.EncodeToString(mac.Sum(nil)))

	// Compare HMAC with received request
	if hmac.Equal(signatureByte, mac.Sum(nil)) {
		fmt.Println("Valid invocation (with validated content) on " + r.URL.Path)
		_, err = fmt.Fprint(w, "Successfully invoked dummy web service.")
		if err != nil {
			fmt.Println("Something went wrong when sending response: " + err.Error())
		}
	} else { // Error - invalid HMAC
		fmt.Println("Invalid invocation (tampered content?) on " + r.URL.Path)
		http.Error(w, "Invalid invocation", http.StatusBadRequest)
	}
}

func main() {


	port := "8081"

	// Environment variable constant for Heroku support
	PORT := "PORT"

	if os.Getenv(PORT) != "" {
		port = os.Getenv(PORT)
	}

	endpoint := "/invoked/"

	fmt.Println("Service listening on port " +  port)
	switch validationLevel {
		case 0:
			fmt.Println("Service URL (non-validating): http://localhost:" + port + endpoint)
			http.HandleFunc(endpoint, NonValidatingHandler)
		case 1:
			fmt.Println("Service URL (URL-validating): http://localhost:" + port + endpoint + hex.EncodeToString(urlMac.Sum(nil)))
			http.HandleFunc(endpoint, URLValidatingHandler)
		case 2:
			fmt.Println("Service URL (content-validating): http://localhost:" + port + endpoint)
			http.HandleFunc(endpoint, ContentValidatingHandler)
		default:
			log.Fatal("Invalid validation level. Exiting ...")

	}
	log.Fatal(http.ListenAndServe(":" + port, nil))

}
