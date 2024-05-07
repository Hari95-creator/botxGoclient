package model

import (
	"database/sql"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"
	"mime/multipart"

	"github.com/google/uuid"
)

type Customer struct {
	ID           int
	GID          string
	PHONE_NUMBER string
	NAME         string
	CREATED_DATE time.Time
}

type Contacts struct {
	COUNTRY_CODE int
	PHONE_NUMBER string
	NAME         string
	EMAIL        string
}

type CSVRepository interface {
	ReadDataFromCSV(filename string) ([]*Contacts, map[string]interface{}, error)
	ReadDataFromCSVFile(file multipart.File) ([]*Contacts, map[string]interface{}, error)
}

type CustomerRepository interface {
	CustomerList(offset, limit int) ([]*Customer, error)
}

type customerRepo struct {
	db *sql.DB
}

type csvRepo struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepo{db: db}
}

func NewCsvRepository(db *sql.DB) CSVRepository {
	return &csvRepo{db: db}
}

func (cu *customerRepo) CustomerList(offset, limit int) ([]*Customer, error) {
	var customers []*Customer

	// Modify the SQL query to include pagination
	query := "SELECT id, gid, phone_number, name, created_date FROM public.customer LIMIT $1 OFFSET $2"
	rows, err := cu.db.Query(query, limit, offset)
	if err != nil {
		log.Println("Error retrieving customers from database:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.GID, &customer.PHONE_NUMBER, &customer.NAME, &customer.CREATED_DATE)
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

func (cu *csvRepo) ReadDataFromCSV(filename string) ([]*Contacts, map[string]interface{}, error) {
	var contacts []*Contacts

	// Open the CSV file
	filePath := "C:/Users/HARI KRISHNAN SG/Desktop/democsv/" + filename
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening CSV file:", err)
		return nil, nil, err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Skip the header row if it exists
	if _, err := reader.Read(); err != nil {
		log.Println("Error reading CSV header:", err)
		return nil, nil, err
	}

	// Prepare the SQL statement for inserting data
	stmt, err := cu.db.Prepare("INSERT INTO public.customer (gid, phone_number, name, created_date, country_code, email) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return nil, nil, err
	}
	defer stmt.Close()

	// Read each record from the CSV file
	for {
		record, err := reader.Read()
		if err == io.EOF {
			// End of file
			break
		}
		if err != nil {
			log.Println("Error reading CSV record:", err)
			return nil, nil, err
		}

		// Check if the record has the expected number of fields
		if len(record) < 4 {
			log.Println("Invalid CSV record:", record)
			continue // Skip this record
		}

		// Generate a random UUID for gid
		gid := uuid.New()

		// Parse the country code
		countryCode, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println("Error converting country code to integer:", err)
			return nil, nil, err
		}

		// Execute the SQL statement to insert data into the table
		_, err = stmt.Exec(gid, record[1], record[2], time.Now(), countryCode, record[3])
		if err != nil {
			log.Println("Error inserting data into the database:", err)
			return nil, nil, err
		}

		// Create a new contact instance and append it to the contacts slice
		contact := &Contacts{
			COUNTRY_CODE: countryCode,
			PHONE_NUMBER: record[1],
			NAME:         record[2],
			EMAIL:        record[3],
		}
		contacts = append(contacts, contact)
	}

	// Create a success response JSON with contacts
	response := map[string]interface{}{
		"status":   "success",
		"message":  "Data inserted successfully",
		"date":     time.Now().Format(time.RFC3339),
		"contacts": contacts,
	}

	return contacts, response, nil
}

func (cu *csvRepo) ReadDataFromCSVFile(file multipart.File) ([]*Contacts, map[string]interface{}, error) {
	var contacts []*Contacts

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Skip the header row if it exists
	if _, err := reader.Read(); err != nil {
		log.Println("Error reading CSV header:", err)
		return nil, nil, err
	}

	// Prepare the SQL statement for inserting data
	stmt, err := cu.db.Prepare("INSERT INTO public.customer (gid, phone_number, name, created_date, country_code, email) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		log.Println("Error preparing SQL statement:", err)
		return nil, nil, err
	}
	defer stmt.Close()

	// Read each record from the CSV file
	for {
		record, err := reader.Read()
		if err == io.EOF {
			// End of file
			break
		}
		if err != nil {
			log.Println("Error reading CSV record:", err)
			return nil, nil, err
		}

		// Check if the record has the expected number of fields
		if len(record) < 4 {
			log.Println("Invalid CSV record:", record)
			continue // Skip this record
		}

		// Generate a random UUID for gid
		gid := uuid.New()

		// Parse the country code
		countryCode, err := strconv.Atoi(record[0])
		if err != nil {
			log.Println("Error converting country code to integer:", err)
			return nil, nil, err
		}

		// Execute the SQL statement to insert data into the table
		_, err = stmt.Exec(gid, record[1], record[2], time.Now(), countryCode, record[3])
		if err != nil {
			log.Println("Error inserting data into the database:", err)
			return nil, nil, err
		}

		// Create a new contact instance and append it to the contacts slice
		contact := &Contacts{
			COUNTRY_CODE: countryCode,
			PHONE_NUMBER: record[1],
			NAME:         record[2],
			EMAIL:        record[3],
		}
		contacts = append(contacts, contact)
	}

	// Create a success response JSON with contacts
	response := map[string]interface{}{
		"status":   "success",
		"message":  "Data inserted successfully",
		"date":     time.Now().Format(time.RFC3339),
		"contacts": contacts,
	}

	return contacts, response, nil
}

