package domain

import (
	"encoding/json"
)

type Product struct {
	Id          int     `json:"id" validate:"required"`
	Name        string  `json:"name"`
	Quantity    *int    `json:"quantity,omitempty"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type ProductResponse struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
	Success     bool    `json:"success"`
}

func (p Product) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(b)
}
