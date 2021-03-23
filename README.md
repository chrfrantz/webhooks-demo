# webhooks-demo

Webhooks demo for teaching purposes


## Basic Usage Instructions

### Version 1 (internal client)

Registration:

- Start server.go to provide webhook registry (localhost:8080/webhook)

- Start client.go (first review whether validation should be activated - see variable "validationLevel")

- Send post request to localhost:8080/webhook with body

{
    "url": "http://localhost:8081/invokedUrl",
    "event": "POST"
}

(Note: Take URL from console output when service is started - which deviates for validating version)

Invocation:

- Invoke localhost:8080/service with POST request

- Observe invocation of registered webhook


### Version 2 (external client)

- Get webhook URL from https://webhook.site/

- Send post request to localhost:8080/webhook with body

{
    "url": "url from webhook.site",
    "event": "POST"
}

- Now invoke as describe in Version 1. The external webhook should be invoked.