package main

import (
	"database/sql"
	"log"
	"net/http"
	"whatbot/controller"
	dbconfig "whatbot/dbConfig"
	"whatbot/model"

	_ "github.com/lib/pq"
)

func main() {
	// Initialize database connection
	db, err := sql.Open("postgres", dbconfig.ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := model.NewUserRepository(db)

	userController := controller.NewUserController(userRepository)

	customerRepository := model.NewCustomerRepository(db)

	customerController := controller.NewCustomerController(customerRepository)

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

	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Server Started In Port 8080")
}
