package controller

import (
	"encoding/json"
	"log"
	"net/http"
	envconfig "whatbot/dbConfig" // This should match the package name and path where envconfig.go is located
	"whatbot/model"
)

type TemplateController struct{}

func (tc *TemplateController) GetAllTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	// Load configuration using the new envconfig package
	config, err := envconfig.LoadConfig("config.json")
	if err != nil {
		log.Panicf("conferr")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	templates, err := model.GetAllTemplates(config.AccessToken, config.WabaId)

	if err != nil {
		http.Error(w, "Failed to fetch templates", http.StatusBadRequest)
		return
	}

	templatesJSON, err := json.Marshal(templates)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(templatesJSON)
}
