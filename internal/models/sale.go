package models

import "time"

type Sale struct {
    ID        string    `json:"id"`
    ProductID string    `json:"product_id"`
    Quantity  int       `json:"quantity"`
    Total     float64   `json:"total"`
    Date      time.Time `json:"date"`
}