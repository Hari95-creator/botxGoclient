package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	envconfig "whatbot/dbConfig"
)

type WhatsAppMessageData struct {
	MessagingProduct string    `json:"messaging_product"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Contact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type Message struct {
	ID            string `json:"id"`
	MessageStatus string `json:"message_status"`
}

func SendMsg(templatename string, recPhone string) (*WhatsAppMessageData, error) {
	config, err := envconfig.LoadConfig("config.json")
	if err != nil {
		return nil, err
	}
	payloadFormat := `{
		"messaging_product": "whatsapp",
		"recipient_type": "individual",
		"to": "%s",
		"type": "template",
		"template": {
			"name": "%s",
			"language": {
				"code": "en_US"
			},
			"components": [
				{
					"type": "body",
					"parameters": []
				}
			]
		}
	}`
	payload := fmt.Sprintf(payloadFormat, recPhone, templatename)
	requestBody := strings.NewReader(payload)
	url := fmt.Sprintf("%s/%s/%s/messages", config.Url, config.Version, config.PhoneNumberId)
	request, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.AccessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to Send Message: %s", string(body))
	}

	var msgsend WhatsAppMessageData
	if err := json.NewDecoder(response.Body).Decode(&msgsend); err != nil {
		return nil, err
	}

	return &msgsend, nil
}
