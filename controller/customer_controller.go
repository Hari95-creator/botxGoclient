package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"whatbot/model"
)

type CustomerController struct {
	CustomerService model.CustomerRepository
}

func NewCustomerController(customerService model.CustomerRepository) *CustomerController {
	return &CustomerController{CustomerService: customerService}
}

func (customer *CustomerController) ListAllCustomer(w http.ResponseWriter, r *http.Request) {

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

	customers, err := customer.CustomerService.CustomerList(offset, pageSize)
	if err != nil {
		log.Println("Error fetching customers:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	customerJSON, err := json.Marshal(customers)
	if err != nil {
		log.Println("Error marshalling customers to JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(customerJSON)
}
