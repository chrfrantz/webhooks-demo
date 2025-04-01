package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"webhooks/handlers"
	structs "webhooks/webhooks"
)

func init() {
	// Secret (relevant for content validation)
	handlers.Secret = []byte{1, 2, 3, 4, 5} // not a good secret!
}

func main() {

	port := "8080"

	// Environment variable constant for PaaS support
	PORT := "PORT"

	if os.Getenv(PORT) != "" {
		port = os.Getenv(PORT)
	}

	// Register default handler
	http.HandleFunc("/", handlers.DefaultServerHandler)
	// Register webhook registration handler
	http.HandleFunc(structs.WebhookEndpoint, handlers.WebhookHandler)
	// Register invocation handler (that triggers invocation of registered webhooks)
	http.HandleFunc(structs.ServiceEndpoint, handlers.ServiceHandler)
	fmt.Println("Service listening on port " + port)
	fmt.Println("For registration of webhook, send POST request to http://localhost:" + port + structs.WebhookEndpoint)
	fmt.Println("For invocation of any registered webhook, send POST request to http://localhost:" + port + structs.ServiceEndpoint)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
