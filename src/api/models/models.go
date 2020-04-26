package models

type TransactionType string

// Define ENUM for constant values "credit" and "debit" operation
const (
	Credit TransactionType = "credit"
	Debit TransactionType = "debit"
)

type Transaction struct {
	ID string `json:"id_transaction"`
	Type TransactionType `json:"type"`
	Ammount float64 `json:"ammount"`
	Date string `json:"date"`
}

type TransactionBody struct {
	Type TransactionType `json:"type"`
	Ammount float64 `json:"ammount"`
}

type AccountBalance struct {
	Balance float64 `json:"balance"`
}