// controller/user_controller.go

package controller

import (
	"encoding/json"
	"net/http"
	"time"
	"whatbot/model"
	"whatbot/utils"

	"github.com/golang-jwt/jwt/v4"
)

type AuthResponse struct {
	Message  string `json:"message"`
	UserID   int    `json:"userId"`
	UserName string `json:"username"`
}

type Claims struct {
	UserID   int    `json:"userId"`
	Username string `json:"username"`
	UserGid  string `json:"gid"`
	jwt.StandardClaims
}

type UserController struct {
	UserService model.UserRepository
}

func NewUserController(userService model.UserRepository) *UserController {
	return &UserController{UserService: userService}
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
		// Generate JWT token
		token := jwt.New(jwt.SigningMethodRS256)
		claims := &Claims{
			UserID:   authenticated.ID,
			Username: authenticated.UserName,
			UserGid:  authenticated.GID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			},
		}
		token.Claims = claims

		// Sign the token with the RSA private key
		tokenString, err := token.SignedString(utils.GetClientPrivateKey())
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
