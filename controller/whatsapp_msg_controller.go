package controller

import (
	"encoding/json"
	"net/http"
	"whatbot/model"
)

func (tc *TemplateController) SendsingleMsg(w http.ResponseWriter, r *http.Request) {

	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	recNumber, ok := requestBody["recNumber"]
	if !ok {
		http.Error(w, "Missing recNumber in request body", http.StatusBadRequest)
		return
	}

	templateName, ok := requestBody["templateName"]
	if !ok {
		http.Error(w, "Missing templatename in request body", http.StatusBadRequest)
		return
	}
	msgsend, err := model.SendMsg(templateName, recNumber)
	if err != nil {
		http.Error(w, "Failed to Send Message", http.StatusBadRequest)
		return
	}
	templatesJSON, err := json.Marshal(msgsend)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(templatesJSON)
}
