// main.go

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

	// Initialize user repository
	userRepository := model.NewUserRepository(db)

	// Initialize user controller
	userController := controller.NewUserController(userRepository)

	// // Define routes
	// http.HandleFunc("/login", userController.Login)

	// // Start the server
	// log.Fatal(http.ListenAndServe(":8080", nil))

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}

	// Define routes
	http.Handle("/login", corsMiddleware(http.HandlerFunc(userController.Login)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
