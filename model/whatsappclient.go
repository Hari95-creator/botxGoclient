package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TemplateData struct {
	Data   []Template `json:"data"`
	Paging PagingInfo `json:"paging"`
}

type Template struct {
	Name       string      `json:"name"`
	Components []Component `json:"components"`
	Language   string      `json:"language"`
	Status     string      `json:"status"`
	Category   string      `json:"category"`
	ID         string      `json:"id"`
}

type Component struct {
	Type    string   `json:"type"`
	Format  string   `json:"format,omitempty"`
	Text    string   `json:"text,omitempty"`
	Example *Example `json:"example,omitempty"`
	Buttons []Button `json:"buttons,omitempty"`
}

type Example struct {
	HeaderHandle []string `json:"header_handle"`
}

type Button struct {
	Type string `json:"type"`
	Text string `json:"text"`
	URL  string `json:"url"`
}

type PagingInfo struct {
	Cursors Cursor `json:"cursors"`
}

type Cursor struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

func GetAllTemplates(accessToken, wabaID string) (*TemplateData, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/message_templates", wabaID)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(response.Body)
		return nil, fmt.Errorf("failed to fetch templates: %s", string(body))
	}

	var templates TemplateData
	if err := json.NewDecoder(response.Body).Decode(&templates); err != nil {
		return nil, err
	}

	return &templates, nil
}
