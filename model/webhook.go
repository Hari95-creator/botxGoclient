package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type webHookRepository interface{

	webhookHandler()
}

type WebhookPayload struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					Type string `json:"type"`
					Text struct {
						Text string `json:"text"`
					} `json:"text,omitempty"`
				} `json:"messages,omitempty"`
				Metadata struct {
					PhoneNumberID string `json:"phone_number_id"`
				} `json:"metadata,omitempty"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Log incoming messages
	fmt.Printf("Incoming webhook message: %+v\n", payload)

	// Check if the webhook request contains a message
	if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 && len(payload.Entry[0].Changes[0].Value.Messages) > 0 {
		message := payload.Entry[0].Changes[0].Value.Messages[0]
		// Check if the incoming message contains text
		if message.Type == "text" {
			// Extract the business number to send the reply from it
			businessPhoneNumberID := payload.Entry[0].Changes[0].Value.Metadata.PhoneNumberID

			// Send a reply message
			// Implement sending reply message logic here

			log.Println(businessPhoneNumberID)
		}
	}

	w.WriteHeader(http.StatusOK)
}
