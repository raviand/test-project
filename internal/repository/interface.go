package repository

import "github.com/raviand/test-project/internal/domain"

type DataInterface interface {
	CreateProduct(product *domain.Product) (*domain.Product, error)
	GetProductById(id int) (*domain.Product, error)
	PatchProduct(id int, product *domain.Product) (*domain.Product, error)
}
