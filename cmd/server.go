package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"webhooks"
)

func init() {
	// Secret (relevant for content validation)
	webhooks.Secret = []byte{1, 2, 3, 4, 5} // not a good secret!
}

func main() {

	port := "8080"

	// Environment variable constant for Heroku support
	PORT := "PORT"

	if os.Getenv(PORT) != "" {
		port = os.Getenv(PORT)
	}

	// for registration
	webhookEndpoint := "/webhook"
	// for invocation
	serviceEndpoint := "/service"

	http.HandleFunc(webhookEndpoint, webhooks.WebhookHandler)
	http.HandleFunc(serviceEndpoint, webhooks.ServiceHandler)
	fmt.Println("Service listening on port " +  port)
	fmt.Println("For registration of webhook, send POST request to http://localhost:" + port + webhookEndpoint)
	fmt.Println("For invocation of any registered webhook, send POST request to http://localhost:" + port + serviceEndpoint)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}



