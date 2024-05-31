package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"net"
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
	CountryService  model.CountryRepository
}

type Customer struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	CreatedDate string `json:"created_date"`
	GID         string `json:"gid"`
}

type Pagination struct {
	Page         int    `json:"page"`
	FirstPageURL string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	Links        []Link `json:"links"`
	NextPageURL  string `json:"next_page_url"`
	ItemsPerPage int    `json:"items_per_page"`
	PrevPageURL  string `json:"prev_page_url"`
	To           int    `json:"to"`
	Total        int    `json:"total"`
}

type Link struct {
	URL    string `json:"url"`
	Label  string `json:"label"`
	Active bool   `json:"active"`
	Page   int    `json:"page,omitempty"`
}

type ResponsePayload struct {
	Data    []Customer `json:"data"`
	Payload struct {
		Pagination Pagination `json:"pagination"`
	} `json:"payload"`
}

func NewCustomerController(customerService model.CustomerRepository, csvService model.CSVRepository, countryService model.CountryRepository) *CustomerController {
	return &CustomerController{
		CustomerService: customerService,
		CSVService:      csvService,
		CountryService:  countryService,
	}
}

func (customer *CustomerController) ListAllCustomer(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	customers, totalCustomers, err := customer.CustomerService.CustomerList(offset, pageSize)
	if err != nil {
		log.Println("Error fetching customers:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var responseCustomers []Customer
	for _, c := range customers {
		responseCustomers = append(responseCustomers, Customer{
			ID:          c.ID,
			Name:        c.NAME,
			PhoneNumber: c.PHONE_NUMBER,
			CreatedDate: c.CREATED_DATE.Format("2006-01-02"),
			GID:         c.GID,
		})
	}
	lastPage := (totalCustomers + pageSize - 1) / pageSize
	var links []Link

	if page > 1 {
		links = append(links, Link{
			URL:    "/?page=" + strconv.Itoa(page-1),
			Label:  "&laquo; Previous",
			Active: false,
			Page:   page - 1,
		})
	}

	for i := 1; i <= lastPage; i++ {
		links = append(links, Link{
			URL:    "/?page=" + strconv.Itoa(i),
			Label:  strconv.Itoa(i),
			Active: i == page,
			Page:   i,
		})
	}

	if page < lastPage {
		links = append(links, Link{
			URL:    "/?page=" + strconv.Itoa(page+1),
			Label:  "Next &raquo;",
			Active: false,
			Page:   page + 1,
		})
	}
	response := ResponsePayload{
		Data: responseCustomers,
		Payload: struct {
			Pagination Pagination `json:"pagination"`
		}{
			Pagination: Pagination{
				Page:         page,
				FirstPageURL: "/?page=1",
				From:         offset + 1,
				LastPage:     lastPage,
				Links:        links,
				NextPageURL:  "/?page=" + strconv.Itoa(page+1),
				ItemsPerPage: pageSize,
				PrevPageURL:  "/?page=" + strconv.Itoa(page-1),
				To:           offset + len(responseCustomers),
				Total:        totalCustomers,
			},
		},
	}
	customerJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling customers to JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(customerJSON)
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
		host, _, err := net.SplitHostPort(ip)
		if err == nil {
			ip = host
		}
	}

	// token  validation part -starts

	// tokenValidationErr := utils.IsTokenValid(tokenString)
	// if tokenValidationErr != nil {
	// 	http.Error(w, tokenValidationErr.Error(), http.StatusUnauthorized)
	// 	return
	// }

	claims := &Claims{}

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
	if claims.ExpiresAt.Time.Before(time.Now()) {
		http.Error(w, "Token has expired", http.StatusUnauthorized)
		return
	}

	if claims.CreationDate.IsZero() || time.Now().Before(claims.CreationDate) {
		http.Error(w, "Invalid token creation date", http.StatusUnauthorized)
		return
	}

	// token validation-ends

	_, successResponse, err := customer.CSVService.ReadDataFromCSVFile(file, claims.UserID, ip)
	if err != nil {
		log.Println("Error reading data from CSV:", err)
		http.Error(w, "Error reading data from CSV", http.StatusInternalServerError)
		return
	}

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

func (customer *CustomerController) CountriesHandler(w http.ResponseWriter, r *http.Request) {
	
	countryCode := r.URL.Query().Get("code")
	if countryCode != "" {
		countries, err := customer.CountryService.GetCountriesByCode(countryCode)
		if err != nil {
			http.Error(w, "Failed to fetch countries", http.StatusInternalServerError)
			return
		}

		if len(countries) == 0 {
			http.Error(w, "No countries found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countries)
		return
	}
	countries, err := customer.CountryService.GetAllCountries()
	if err != nil {
		http.Error(w, "Failed to fetch countries", http.StatusInternalServerError)
		return
	}

	if countries == nil {
		http.Error(w, "No countries found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(countries)
}

