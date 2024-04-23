package model

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type TemplateData struct {
    Data []struct {
        Name       string   `json:"name"`
        Language   string   `json:"language"`
        Category   string   `json:"category"`
        Components []string `json:"components"`
    } `json:"data"`
}

func GetAllTemplates(accessToken, wabaID string) (*TemplateData, error) {
    url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/message_templates", wabaID)
    request, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    request.Header.Set("Authorization", "Bearer "+accessToken)
    request.Header.Set("Content-Type", "application/json")

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
