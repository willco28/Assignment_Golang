package entity

import "time"

type SaldoWallet struct {
	ID     int     `json:"id"`
	IDUser int     `json:"id_user" binding:"required"`
	Name   string  `json:"name" binding:"required,min=3"`
	Saldo  float64 `json:"saldo"`
}

type HistoryTransaction struct {
	ID              int       `json:"id"`
	IDUser          int       `json:"id_user" binding:"required"`
	Name            string    `json:"name" binding:"required,min=3"`
	Saldo           float64   `json:"saldo"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	TransactionDate time.Time `json:"Transaction_date"`
}
