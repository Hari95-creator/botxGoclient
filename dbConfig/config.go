package dbconfig

import (
	"fmt"
)

var (
	host     string
	port     int
	user     string
	password string
	dbname   string
)

func init() {

	host = "localhost"
	port = 5432
	user = "postgres"
	password = "postgres"
	dbname = "WhatBot"
}

func ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
