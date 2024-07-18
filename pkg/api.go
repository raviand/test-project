package pkg

import "time"

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CodeValue   string    `json:"code_value"`
	IsPublished bool      `json:"is_published"`
	Expiration  time.Time `json:"expiration"`
}

type ProductPatchRequest struct {
	ID          *int     `json:"id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Quantity    *int     `json:"quantity,omitempty"`
	CodeValue   *string  `json:"code_value,omitempty"`
	IsPublished *bool    `json:"is_published,omitempty"`
	Expiration  *string  `json:"expiration,omitempty"`
}
