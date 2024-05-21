package models

import (
	"time"
)


type DBConfig struct {
	DB_User     string `json:"db_user"`
	DB_Pass     string `json:"db_pass"`
	DB_Host     string `json:"db_host"`
	DB_Port     string `json:"db_port"`
	DB_Database string `json:"db_database"`
	DB_Sslmode  string `json:"db_sslmode"`
}

type Order struct {
    OrderID    int       `json:"order_id"`   
    UserID     int       `json:"user_id"`     
    ProductID  int       `json:"product_id"`  
    Quantity   int       `json:"quantity"`    
    TotalPrice float64   `json:"total_price"` 
    Status     string    `json:"status"`      // Current status of the order (e.g., pending, completed, cancelled)
    CreatedAt  time.Time `json:"created_at"`  
}


type Payment struct {
	PaymentID        int       `json:"payment_id"`
	OrderID          int       `json:"order_id"`
	Amount           float64   `json:"amount"`
	Status           string    `json:"status"`
	TransactionDate  time.Time `json:"transaction_date"`
}


type JsonResponse struct {
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"status_code"`
}

type ErrorResponse struct {
	Status     bool   `json:"status"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
