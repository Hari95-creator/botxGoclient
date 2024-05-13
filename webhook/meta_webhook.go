package webhook

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// HandlerFunc is a function type for handling webhook requests
type HandlerFunc func(http.ResponseWriter, *http.Request, *sql.DB, WebhookPayload)

type WebhookPayload struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type     string `json:"type"`
					Metadata struct {
						DisplayPhoneNumber string `json:"display_phone_number"`
						PhoneNumberID      string `json:"phone_number_id"`
					} `json:"metadata"`
				} `json:"messages"`
			} `json:"value"`
			Field string `json:"field"`
		} `json:"changes"`
	} `json:"entry"`
}

// WebhookHandler creates a handler function for webhook requests
func WebhookHandler(handler HandlerFunc, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload WebhookPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler(w, r, db, payload)
	}
}

func InsertWhatsappMsgData(db *sql.DB, jsonData []byte) error {
	gid := uuid.New()

	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	messages := data["messages"].([]interface{})
	message := messages[0].(map[string]interface{})
	from := message["from"].(string)

	metadata := data["metadata"].([]interface{})
	bussinessId := metadata[0].(map[string]interface{})["bussinessId"].(string)
	phoneNumberID := metadata[0].(map[string]interface{})["phoneNumberID"].(string)

	// take when need (will apply based on logic)
	// msgBody := message["msgBody"].(string)
	// timestamp := message["timestamp"].(string)

	// // Extracting data from contacts and metadata
	// contacts := data["contacts"].([]interface{})

	// _, err = db.Exec("INSERT INTO whatsapp_data (gid, sender_phone_number, message_body, message_timestamp, display_phone_number) VALUES ($1, $2, $3, $4, $5)", gid, from, msgBody, timestamp, phoneNumber)
	// return err

	_, err = db.Exec("INSERT INTO whatsapp_data (gid,bussiness_id,phone_number_id,sender_phone_number, message_data) VALUES ($1,$2,$3,$4,$5)", gid, bussinessId, phoneNumberID, from, string(jsonData))
	return err
}
