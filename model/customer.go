package model

import (
    "database/sql"
    "log"
    "time"
)

type Customer struct {
    ID           int
    GID          string
    PHONE_NUMBER string
	NAME         string
    CREATED_DATE time.Time
}

type CustomerRepository interface {
    CustomerList() ([]*Customer, error)
}

type customerRepo struct {
    db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
    return &customerRepo{db: db}
}

func (cu *customerRepo) CustomerList() ([]*Customer, error) {
    var customers []*Customer

    rows, err := cu.db.Query("SELECT id, gid ,phone_number,name,created_date FROM public.customer")
    if err != nil {
        log.Println("Error retrieving customers from database:", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var customer Customer
        err := rows.Scan(&customer.ID, &customer.GID, &customer.PHONE_NUMBER,&customer.NAME, &customer.CREATED_DATE)
        if err != nil {
            log.Println("Error scanning customer row:", err)
            continue
        }
        customers = append(customers, &customer)
    }

    if err := rows.Err(); err != nil {
        log.Println("Error iterating over customer rows:", err)
        return nil, err
    }

    return customers, nil
}
