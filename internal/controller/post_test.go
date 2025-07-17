package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/raviand/test-project/pkg"
	"github.com/raviand/test-project/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	mockService := new(test.MockProductService)
	handler := NewHandler(mockService, nil) // audit channel is nil for this test

	tests := []struct {
		Title   string
		Payload *http.Request
		Error   error
		Check   func(response *http.Response)
	}{
		{
			Title: "Post Test",
			Payload: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString(`{}`)),
			},
			Error: nil,
			Check: func(response *http.Response) {
				assert.Equal(t, 200, response.StatusCode)
			},
		},
		{
			Title: "Post Test 2",
			Payload: &http.Request{
				Body: io.NopCloser(bytes.NewBufferString(`{}`)),
			}, Error: nil,
			Check: func(response *http.Response) {
				assert.Equal(t, 400, response.StatusCode)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Title, func(t *testing.T) {
			res := httptest.NewRecorder()
			handler.Create(res, test.Payload)
			test.Check(res.Result())
		})
	}

	t.Run("Success", func(t *testing.T) {
		productReq := pkg.CreateProductRequest{
			Name:        "Test Product",
			Price:       10.0,
			Quantity:    100,
			CodeValue:   "TP001",
			IsPublished: true,
			Expiration:  "01/01/2025",
			ProductType: 1,
		}

		parsedTime, _ := time.Parse("02/01/2006", productReq.Expiration)
		expectedProduct := &pkg.Product{
			Name:          productReq.Name,
			Price:         productReq.Price,
			Quantity:      productReq.Quantity,
			CodeValue:     productReq.CodeValue,
			IsPublished:   productReq.IsPublished,
			Expiration:    parsedTime,
			ProductTypeId: productReq.ProductType,
		}

		mockService.On("Create", expectedProduct).Return(nil).Once()

		body, _ := json.Marshal(productReq)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("{invalid"))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var errResp pkg.ApiError
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Equal(t, pkg.ErrorCode(pkg.BadRequest), errResp.Code)
	})

	t.Run("Bad Request - Invalid Date Format", func(t *testing.T) {
		productReq := pkg.CreateProductRequest{
			Expiration: "2025-01-01", // wrong format
		}
		body, _ := json.Marshal(productReq)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var errResp pkg.ApiError
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Equal(t, pkg.ErrorCode(pkg.WrongFieldValues), errResp.Code)
	})

	t.Run("Internal Server Error on Create", func(t *testing.T) {
		productReq := pkg.CreateProductRequest{
			Name:        "Test Product",
			Price:       10.0,
			Quantity:    100,
			CodeValue:   "TP001",
			IsPublished: true,
			Expiration:  "01/01/2025",
			ProductType: 1,
		}

		parsedTime, _ := time.Parse("02/01/2006", productReq.Expiration)
		expectedProduct := &pkg.Product{
			Name:          productReq.Name,
			Price:         productReq.Price,
			Quantity:      productReq.Quantity,
			CodeValue:     productReq.CodeValue,
			IsPublished:   productReq.IsPublished,
			Expiration:    parsedTime,
			ProductTypeId: productReq.ProductType,
		}

		expectedError := pkg.GetError(pkg.InternalError)
		mockService.On("Create", expectedProduct).Return(expectedError).Once()

		body, _ := json.Marshal(productReq)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var errResp pkg.ApiError
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Equal(t, pkg.ErrorCode(pkg.InternalError), errResp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestCreateUser(t *testing.T) {
	mockService := new(test.MockProductService)
	handler := NewHandler(mockService, nil)

	t.Run("Success", func(t *testing.T) {
		user := pkg.User{
			Id:        "1",
			Firstname: "John",
			Lastname:  "Doe",
		}
		mockService.On("CreateUser", mock.Anything, &user).Return(nil).Once()

		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.CreateUser(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("{invalid"))
		rr := httptest.NewRecorder()

		handler.CreateUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var errResp pkg.ApiError
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Equal(t, pkg.ErrorCode(pkg.BadRequest), errResp.Code)
	})

	t.Run("Internal Server Error on CreateUser", func(t *testing.T) {
		user := pkg.User{
			Id:        "1",
			Firstname: "John",
			Lastname:  "Doe",
		}
		expectedError := pkg.GetError(pkg.InternalError)
		mockService.On("CreateUser", mock.Anything, &user).Return(expectedError).Once()

		body, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler.CreateUser(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var errResp pkg.ApiError
		json.Unmarshal(rr.Body.Bytes(), &errResp)
		assert.Equal(t, pkg.ErrorCode(pkg.InternalError), errResp.Code)
		mockService.AssertExpectations(t)
	})
}
