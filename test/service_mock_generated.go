package test

import (
	"context"

	"github.com/raviand/test-project/pkg"
	"github.com/stretchr/testify/mock"
)

// MockProductService is a mock type for the ProductService type
type MockProductService struct {
	mock.Mock
}

// Create provides a mock function with given fields: product
func (m *MockProductService) Create(product *pkg.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

// CreateUser provides a mock function with given fields: ctx, user
func (m *MockProductService) CreateUser(ctx context.Context, user *pkg.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Delete provides a mock function with given fields: id
func (m *MockProductService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// DeleteUser provides a mock function with given fields: ctx, id
func (m *MockProductService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetAll provides a mock function with given fields:
func (m *MockProductService) GetAll() ([]*pkg.Product, error) {
	args := m.Called()
	return args.Get(0).([]*pkg.Product), args.Error(1)
}

// GetByID provides a mock function with given fields: id
func (m *MockProductService) GetByID(id int) (*pkg.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*pkg.Product), args.Error(1)
}

// GetUser provides a mock function with given fields: ctx, id
func (m *MockProductService) GetUser(ctx context.Context, id string) (*pkg.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*pkg.User), args.Error(1)
}

// ListAllUsers provides a mock function with given fields: ctx
func (m *MockProductService) ListAllUsers(ctx context.Context) ([]pkg.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]pkg.User), args.Error(1)
}

// Patch provides a mock function with given fields: id, product
func (m *MockProductService) Patch(id int, product *pkg.ProductPatchRequest) error {
	args := m.Called(id, product)
	return args.Error(0)
}

// UpdateCreate provides a mock function with given fields: product
func (m *MockProductService) UpdateCreate(product *pkg.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
