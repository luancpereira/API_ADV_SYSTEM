package request

import "time"

/*****
struct for posts
******/

type InsertTransaction struct {
	Description      string    `json:"description"`
	TransactionDate  time.Time `json:"transaction_date"`
	TransactionValue float64   `json:"transaction_value"`
}

/*****
struct for posts
******/
