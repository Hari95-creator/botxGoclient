package controller

import (
    "encoding/json"
    "net/http"
    "whatbot/model"
    "log"
)

type CustomerController struct {
    CustomerService model.CustomerRepository
}

func NewCustomerController(customerService model.CustomerRepository) *CustomerController {
    return &CustomerController{CustomerService: customerService}
}

func (customer *CustomerController) ListAllCustomer(w http.ResponseWriter, r *http.Request) {
    customers, err := customer.CustomerService.CustomerList()
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
