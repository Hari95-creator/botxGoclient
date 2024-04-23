package controller

import (
    "encoding/json"
    "log"
    "net/http"
    "whatbot/model"
	"os"
)

type TemplateController struct{}

func (tc *TemplateController) GetAllTemplatesHandler(w http.ResponseWriter, r *http.Request) {
    
	_, err := os.Stat("config.json")
	if os.IsNotExist(err) {
		log.Println("Config file not found")
		http.Error(w, "Config file not found", http.StatusInternalServerError)
		return
	}

	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	type Config struct {
		AccessToken string `json:"access_token"`
		WabaID      string `json:"waba_id"`
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Printf("Error unmarshalling config data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	accessToken := config.AccessToken
	wabaID := config.WabaID

    templates, err := model.GetAllTemplates(accessToken, wabaID)
    if err != nil {
        log.Printf("Error fetching templates: %v", err)
        http.Error(w, "Failed to fetch templates", http.StatusInternalServerError)
        return
    }

    templatesJSON, err := json.Marshal(templates)
    if err != nil {
        log.Printf("Error marshalling templates to JSON: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(templatesJSON)
}
