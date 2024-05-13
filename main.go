package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"whatbot/controller"
	dbconfig "whatbot/dbConfig"
	"whatbot/model"
	"whatbot/webhook"

	_ "github.com/lib/pq"
)

const (
	PORT = 8080
)

func main() {
	// Initialize database connection
	db, err := sql.Open("postgres", dbconfig.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Start the webhook server
	go StartWebhookServer(db)

	// Start the HTTP server
	StartHTTPServer(db)
}

func webhookHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, payload webhook.WebhookPayload) {
	if payload.Object != "" {
		var messageResponses []map[string]interface{}
		var metaDataResponses []map[string]interface{}
		var contactResponses []map[string]interface{}

		for _, entry := range payload.Entry {
			for _, change := range entry.Changes {
				if len(change.Value.Messages) > 0 {
					BussinessId := entry.ID
					phoneNumberID := change.Value.Metadata.PhoneNumberID
					DisplayPhoneNumber := change.Value.Metadata.DisplayPhoneNumber
					from := change.Value.Messages[0].From
					msgBody := change.Value.Messages[0].Text.Body
					id := change.Value.Messages[0].ID
					timeStamp := change.Value.Messages[0].Timestamp

					messageResponse := map[string]interface{}{
						"from":      from,
						"msgBody":   msgBody,
						"id":        id,
						"timestamp": timeStamp,
					}
					messageResponses = append(messageResponses, messageResponse)

					metaDataResponse := map[string]interface{}{
						"phoneNumberID":      phoneNumberID,
						"displayPhoneNumber": DisplayPhoneNumber,
						"bussinessId":        BussinessId,
					}

					metaDataResponses = append(metaDataResponses, metaDataResponse)
				}

				if len(change.Value.Contacts) > 0 {
					for _, contact := range change.Value.Contacts {
						profileName := contact.Profile.Name
						waID := contact.WaID

						contactResponse := map[string]interface{}{
							"profileName": profileName,
							"waID":        waID,
						}
						contactResponses = append(contactResponses, contactResponse)
					}
				}
			}
		}

		jsonResponse, err := json.Marshal(map[string]interface{}{
			"metadata":   metaDataResponses,
			"messages":   messageResponses,
			"contacts":   contactResponses,
			"status":     "Success",
			"statuscode": http.StatusOK,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = webhook.InsertWhatsappMsgData(db, jsonResponse)
		if err != nil {
			fmt.Println("Error inserting data:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		fmt.Println("Invalid payload")
		w.WriteHeader(http.StatusNotFound)
	}
}
func StartWebhookServer(db *sql.DB) {
	http.HandleFunc("/webhook", webhook.WebhookHandler(webhookHandler, db))
	log.Println("Webhook server started")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func StartHTTPServer(db *sql.DB) {
	userRepository := model.NewUserRepository(db)
	userController := controller.NewUserController(userRepository)

	customerRepository := model.NewCustomerRepository(db)
	csvRepository := model.NewCsvRepository(db)
	customerController := controller.NewCustomerController(customerRepository, csvRepository)

	whatsappController := controller.TemplateController{}

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	http.Handle("/login", corsMiddleware(http.HandlerFunc(userController.Login)))
	http.Handle("/customer/list", corsMiddleware(http.HandlerFunc(customerController.ListAllCustomer)))
	http.Handle("/templates/", corsMiddleware(http.HandlerFunc(whatsappController.GetAllTemplatesHandler)))
	http.Handle("/sendmessage/", corsMiddleware(http.HandlerFunc(whatsappController.SendsingleMsg)))
	http.Handle("/customer/data/csv/", corsMiddleware(http.HandlerFunc(customerController.ReadCsv)))

	log.Printf("Starting HTTP server on port %d...\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
