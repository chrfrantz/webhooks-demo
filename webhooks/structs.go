package structs

type WebhookRegistration struct {
	Url   string `json:"url"`
	Event string `json:"event"`
}
