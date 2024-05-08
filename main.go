package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"whatbot/controller"
	dbconfig "whatbot/dbConfig"
	"whatbot/model"
	"whatbot/webhook"
	"encoding/json"

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

func webhookHandler(w http.ResponseWriter, r *http.Request, payload webhook.WebhookPayload) {
	if payload.Object != "" {
		if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 && len(payload.Entry[0].Changes[0].Value.Messages) > 0 {
			phoneNumberID := payload.Entry[0].Changes[0].Value.Messages[0].Metadata.PhoneNumberID
			from := payload.Entry[0].Changes[0].Value.Messages[0].From
			msgBody := payload.Entry[0].Changes[0].Value.Messages[0].Text.Body

			fmt.Println("Phone number:", phoneNumberID)
			fmt.Println("From:", from)
			fmt.Println("Body:", msgBody)

			response := map[string]interface{}{
				"phoneNumberID": phoneNumberID,
				"from":          from,
				"msgBody":       msgBody,
			}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
			return
			

			// Implement your logic to send a response back
			// Note: You need to implement the functionality to send a message back to the sender.
			// For this example, we're just sending a 200 OK response.
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	fmt.Println("Invalid payload")
	w.WriteHeader(http.StatusNotFound)
}

func StartWebhookServer(db *sql.DB) {
	http.HandleFunc("/webhook", webhook.WebhookHandler(webhookHandler))
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
