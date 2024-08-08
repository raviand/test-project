package service

import (
	"context"
	"time"

	"github.com/raviand/test-project/internal/repository"
	"github.com/raviand/test-project/pkg"
)

type ProductService interface {
	GetAll() ([]*pkg.Product, error)
	GetByID(id int) (*pkg.Product, error)
	Create(product *pkg.Product) error
	UpdateCreate(product *pkg.Product) error
	Patch(id int, product *pkg.ProductPatchRequest) error
	Delete(id int) error
	CreateUser(ctx context.Context, user *pkg.User) error
	GetUser(ctx context.Context, id string) (*pkg.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type productService struct {
	db       repository.Database
	dynamoDB repository.RepositoryDynamo
}

func NewProductService(db repository.Database, dynamo repository.RepositoryDynamo) ProductService {
	return &productService{
		db:       db,
		dynamoDB: dynamo,
	}
}

func (s *productService) GetAll() ([]*pkg.Product, error) {
	return s.db.ListAll()
}

func (s *productService) GetByID(id int) (*pkg.Product, error) {
	if id < 1 {
		return nil, pkg.GetError(pkg.BadRequest)
	}
	return s.db.GetByID(id)
}

func (s *productService) Create(product *pkg.Product) error {
	if product == nil {
		return pkg.GetError(pkg.BadRequest)
	}
	if product.Price < 0 || product.Quantity < 0 {
		return pkg.GetError(pkg.WrongFieldValues)
	}
	return s.db.Create(product)
}

func (s *productService) UpdateCreate(product *pkg.Product) error {
	if product == nil {
		return pkg.GetError(pkg.BadRequest)
	}
	if product.Price < 0 || product.Quantity < 0 {
		return pkg.GetError(pkg.WrongFieldValues)
	}
	if product.ID == 0 {
		return s.db.Create(product)
	}
	_, err := s.db.GetByID(product.ID)
	if err != nil {
		return s.db.Create(product)
	}
	return s.db.Update(product)
}

func (s *productService) Patch(id int, product *pkg.ProductPatchRequest) error {
	if product == nil {
		return pkg.GetError(pkg.BadRequest)
	}
	if product.Price != nil && *product.Price < 0 {
		return pkg.GetError(pkg.WrongFieldValues)
	}
	if product.Quantity != nil && *product.Quantity < 0 {
		return pkg.GetError(pkg.WrongFieldValues)
	}
	actualProduct, err := s.db.GetByID(id)
	if err != nil {
		return err
	}
	newProduct, err := BuildPatch(actualProduct, product)
	if err != nil {
		return err
	}
	return s.db.Update(newProduct)
}

func BuildPatch(product *pkg.Product, patch *pkg.ProductPatchRequest) (*pkg.Product, error) {
	if patch.Name != nil {
		product.Name = *patch.Name
	}
	if patch.Price != nil {
		product.Price = *patch.Price
	}
	if patch.Quantity != nil {
		product.Quantity = *patch.Quantity
	}
	if patch.CodeValue != nil {
		product.CodeValue = *patch.CodeValue
	}
	if patch.IsPublished != nil {
		product.IsPublished = *patch.IsPublished
	}
	if patch.Expiration != nil {
		t, err := time.Parse("02/01/2006", *patch.Expiration)
		if err != nil {
			return nil, pkg.GetError(pkg.WrongFieldValues)
		}
		product.Expiration = t
	}
	return product, nil
}

func (s *productService) Delete(id int) error {
	if id < 1 {
		return pkg.GetError(pkg.BadRequest)
	}
	return s.db.Delete(id)
}

func (s *productService) CreateUser(ctx context.Context, user *pkg.User) error {
	return s.dynamoDB.Store(ctx, user)
}

func (s *productService) GetUser(ctx context.Context, id string) (*pkg.User, error) {
	return s.dynamoDB.GetOne(ctx, id)
}

func (s *productService) DeleteUser(ctx context.Context, id string) error {
	return s.dynamoDB.Delete(ctx, id)
}
