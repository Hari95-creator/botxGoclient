package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"whatbot/model"
	"whatbot/utils"

	"github.com/golang-jwt/jwt/v4"
)

// CustomerController handles HTTP requests related to customers
type CustomerController struct {
	CustomerService model.CustomerRepository
	CSVService      model.CSVRepository
}

// NewCustomerController creates a new instance of CustomerController
func NewCustomerController(customerService model.CustomerRepository, csvService model.CSVRepository) *CustomerController {
	return &CustomerController{CustomerService: customerService, CSVService: csvService}
}

// ListAllCustomer returns a list of customers with pagination
func (customer *CustomerController) ListAllCustomer(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default page if invalid or not provided
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10 // Default page size if invalid or not provided
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Fetch customers from the service layer
	customers, err := customer.CustomerService.CustomerList(offset, pageSize)
	if err != nil {
		log.Println("Error fetching customers:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert customers to JSON
	customerJSON, err := json.Marshal(customers)
	if err != nil {
		log.Println("Error marshalling customers to JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(customerJSON)
}

func (customer *CustomerController) csvFromFile(w http.ResponseWriter, r *http.Request) {

	// Parse the request body to get the CSV file path
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve filename from the request body
	filename, ok := requestBody["filename"]
	if !ok {
		http.Error(w, "Missing filename in request body", http.StatusBadRequest)
		return
	}

	// Read data from CSV file using CSVService
	customers, successResponse, err := customer.CSVService.ReadDataFromCSV(filename)
	if err != nil {
		log.Println("Error reading data from CSV:", err)
		http.Error(w, "Error reading data from CSV", http.StatusInternalServerError)
		return
	}

	log.Println(customers)

	// Create a response map
	response := map[string]interface{}{
		"status":  "success",
		"message": "Data retrieved successfully",
		"data": map[string]interface{}{
			"success": successResponse,
		},
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (customer *CustomerController) ReadCsv(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(30 << 20) // Set a limit on the maximum upload size (30 MB in this example)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tokenString := r.FormValue("token")
	if tokenString == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
	}

	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		ip = parts[0]
	}

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.GetClientPublicKey(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Invalid token signature", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to parse token", http.StatusBadRequest)
		return
	}
	if !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if token is expired
	if claims.ExpiresAt < time.Now().Unix() {
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		return
	}
	// Read data from CSV file using CSVService
	customers, successResponse, err := customer.CSVService.ReadDataFromCSVFile(file, claims.UserID, ip)
	if err != nil {
		log.Println("Error reading data from CSV:", err)
		http.Error(w, "Error reading data from CSV", http.StatusInternalServerError)
		return
	}

	log.Println(customers)
	// Create a response

	// Create a response map
	response := map[string]interface{}{
		"status":  "success",
		"message": "Data retrieved successfully",
		"data": map[string]interface{}{
			"success": successResponse,
		},
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
