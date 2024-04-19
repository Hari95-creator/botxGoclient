// controller/user_controller.go

package controller

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"
	"whatbot/model"

	"encoding/base64"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type AuthResponse struct {
	Message  string `json:"message"`
	UserID   int    `json:"userId"`
	UserName string `json:"username"`
}

type Claims struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type UserController struct {
	UserService model.UserRepository
}

func NewUserController(userService model.UserRepository) *UserController {
	return &UserController{UserService: userService}
}

func generateSecretKey(length int) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode random bytes to base64 string
	key := base64.URLEncoding.EncodeToString(randomBytes)
	return key, nil
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve username and password from the request body
	username, ok := requestBody["username"]
	if !ok {
		http.Error(w, "Missing username in request body", http.StatusBadRequest)
		return
	}
	password, ok := requestBody["password"]
	if !ok {
		http.Error(w, "Missing password in request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	authenticated, err := uc.UserService.AuthenticateUser(username, password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if authenticated != nil {

		// response := AuthResponse{
		// 	Message:  "Login successful",
		// 	UserID:   authenticated.ID,
		// 	UserName: authenticated.UserName,
		// }

		// // Encode the response as JSON
		// jsonResponse, err := json.Marshal(response)
		// if err != nil {
		// 	// Handle JSON encoding error
		// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// 	return
		// }
		// // Authentication successful
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK)
		// w.Write(jsonResponse)

		keyLength := 32
		secretKey, err := generateSecretKey(keyLength)
		if err != nil {
			log.Println("Error generating secret key:", err)
			return
		}
		// Generate JWT token
		token := jwt.New(jwt.GetSigningMethod("HS256"))
		claims := &Claims{
			UserID:   authenticated.ID,
			Username: authenticated.UserName,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
				// ExpiresAt: time.Now().Add(1 * time.Second).Unix(),

			},
		}
		token.Claims = claims

		// Sign the token with a secret key
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the JWT token in the response
		response := map[string]string{
			"token": tokenString,
		}

		// Convert the response to JSON format
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Write the JSON response with the JWT token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	} else {
		// Authentication failed
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}
