package response

import (
	"time"
)

/*****
struct for posts
******/

type ReturnInsertTransaction struct {
	ID int64 `json:"id"`
}

/*****
struct for posts
******/

/*****
struct for gets
******/

type GetTransactions struct {
	ID                                      int64     `json:"id"`
	Description                             string    `json:"description"`
	TransactionDate                         time.Time `json:"transaction_date"`
	TransactionValue                        float64   `json:"transaction_value"`
	ExchangeRate                            float64   `json:"exchange_rate"`
	TransactionValueConvertedToWishCurrency float64   `json:"transaction_value_converted_to_wish_currency"`
}

type GetTransactionsByID struct {
	ID                                      int64     `json:"id"`
	Description                             string    `json:"description"`
	TransactionDate                         time.Time `json:"transaction_date"`
	TransactionValue                        float64   `json:"transaction_value"`
	ExchangeRate                            float64   `json:"exchange_rate"`
	TransactionValueConvertedToWishCurrency float64   `json:"transaction_value_converted_to_wish_currency"`
}

/*****
struct for gets
******/
