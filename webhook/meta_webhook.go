package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HandlerFunc is a function type for handling webhook requests
type HandlerFunc func(http.ResponseWriter, *http.Request, WebhookPayload)

// WebhookPayload represents the structure of the webhook payload
type WebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					From string `json:"from"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
					Metadata struct {
						PhoneNumberID string `json:"phone_number_id"`
					} `json:"metadata"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

// WebhookHandler creates a handler function for webhook requests
func WebhookHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload WebhookPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler(w, r, payload)
	}
}
