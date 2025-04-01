# Webhooks demo

This demo showcases the principles and practical use of webhooks using Golang. This includes the basic publish-subscribe pattern. It further includes optional provisions of URL-based and content-based validation.

## Deployment

### Registration

* Start `server.go` to provide webhook registry (`http://localhost:8080/registration`)

* Start `client.go` (first review whether validation should be activated - see variable "validationLevel")
  * Alternatively, get a webhook URL from https://webhook.site and use that for testing

* To register webhook for the invocation of the client, send `POST` request to `http://localhost:8080/registration` with body

```
{
    "url": "http://localhost:8081/pathToBeInvoked",
    "event": "EVENT"
}
```

  * The content of the `event` field is an event defined in the context of your application the webhook is used in (e.g., updating internal information).

* Test the successful registration by issuing a `GET` request to `http://localhost:8080/registration`.

(Note: Take URL from console output when service is started - which deviates for validating version)

### Invocation

* Invoke `http://localhost:8080/invocation` with `POST` request and arbitrary content (can be plain text, JSON-formatted, etc.)

* Observe invocation of registered webhook (depending on the choice, either in `client.go` or under webhook.site).

### Exploration

Use external URLs to showcase the flexible use of webhooks, such as the following example:

* Get webhook URL from https://webhook.site/

* Register external webhook by sending `POST` request to `http://localhost:8080/registration` with body

```
{
    "url": "url from webhook.site",
    "event": "EXAMPLE EVENT"
}
```

* Now invoke as described under `Invocation`.

Beyond this, explore the inner workings of validations. 

### Advanced topic: Validation

Validation can occur based on URLs (`validationLevel = 1`), as well content (`validationLevel = 2`).

#### URL-based validation
In the first case (`validationLevel = 1`), hashing is used to generate an URL suffix that reflects security by obscurity principles in creating sufficient long URLs that make guessing challenging or impossible. Here, the generated URL (printed on the console during instantiation of the client) needs to be registered on the server webhook registration endpoint. 
* Note that this particular example uses hashes to be able to regenerate hashes based on a shared secret. An alternative implementation of URL-based security could be based on UUIDs. 

#### Content-based validation
The second case (`validationLevel = 2`) hashes the entire content. Here the essential aspect is that the content may well be visible by third parties, but the hashing can provide assurance of information integrity (i.e., that it hasn't been changed in the process of transmission).

## Exercises:
- Run with multiple webhooks (e.g., local and remote ones). Hint: To run multiple local webhooks, compile the client, and run with parameters `port` and `delay` (e.g., `client 8085 0` for running it on port 8085, but without any delay in execution upon message receipt). Ensure to choose different ports if running multiple clients to prevent conflicts (one port can only be bound by one instance).
- Implement mechanism for deletion of webhook
- Create unique identifiers for registered webhooks
- Explore the purpose of goroutines when calling CallUrl (from the server handler) when invoking blocking (e.g., sleeping) clients.
