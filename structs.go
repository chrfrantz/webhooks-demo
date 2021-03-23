package webhooks

type WebhookRegistration struct {
	Url		string	`json:"url"`
	Event 	string	`json:"event"`
}

/*type WebhookInvocation struct {
	Content string `json:"content"`
}*/
