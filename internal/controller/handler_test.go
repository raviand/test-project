package controller

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/raviand/test-project/internal/audit"
	"github.com/raviand/test-project/internal/repository"
	"github.com/raviand/test-project/internal/service"
	"github.com/raviand/test-project/pkg"
	"github.com/raviand/test-project/test"
	"github.com/stretchr/testify/require"
)

func GetErrorReceivedAndExpected(t *testing.T, code pkg.ErrorCode, res *httptest.ResponseRecorder) (pkg.ApiError, pkg.ApiError) {
	errorResponse := pkg.ApiError{}
	require.Nil(t, json.Unmarshal(res.Body.Bytes(), &errorResponse))
	expectedError := pkg.GetError(code)
	return errorResponse, expectedError
}

func TestMiddleware(t *testing.T) {
	// given
	os.Setenv("TOKEN", "this-is-a-token")
	DB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer DB.Close()
	db := repository.NewDatabase(DB)
	svc := service.NewProductService(db)
	audit := make(chan audit.AuditLog)
	handler := NewHandler(svc, audit)
	r := chi.NewRouter()
	r.Use(handler.TokenMiddleware)
	r.Get("/product", handler.GetAll)
	req := httptest.NewRequest("GET", "/product", nil)

	t.Run("should pass the token middleware and succeed the authentication", func(t *testing.T) {
		mock.ExpectQuery("select p.id, p.name, p.quantity, p.price, p.code_value,  p.is_published from product p").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "quantity", "price", "code_value", "is_published"}).AddRow(1, "test", 1, 1, "test", true))
		req.Header.Set("Authorization", "this-is-a-token")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
	})

	t.Run("should reject the request due to token difference", func(t *testing.T) {
		req.Header.Set("Authorization", "this-is-a-different-token")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		require.Equal(t, 401, res.Code)
	})
	close(audit)
}

func TestProductSave(t *testing.T) {
	// given
	DB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("Error initializing mock database: %v", err)
	}
	defer DB.Close()

	db := repository.NewDatabase(DB)
	svc := service.NewProductService(db)
	audit := make(chan audit.AuditLog)
	handler := NewHandler(svc, audit)
	t.Run("should create a new product", func(t *testing.T) {
		p := pkg.CreateProductRequest{
			Name:        "Tv Samsung",
			Price:       1355.8,
			Quantity:    2,
			CodeValue:   "TV_EXT_SAMSUNG",
			IsPublished: false,
			Expiration:  "12/05/2023",
		}
		mock.ExpectExec("INSERT INTO product (name, price, quantity, code_value, is_published, expiration) VALUES (?, ?, ?, ?, ?, ?)").WithArgs(p.Name, p.Price, p.Quantity, p.CodeValue, p.IsPublished, time.Date(2023, 5, 12, 0, 0, 0, 0, time.UTC)).WillReturnResult(sqlmock.NewResult(1, 1))
		b, err := json.Marshal(p)
		require.NoError(t, err)
		req := httptest.NewRequest("POST", "/product", strings.NewReader(string(b)))
		res := httptest.NewRecorder()
		handler.Create(res, req)
		// then
		tim, err := time.Parse("02/01/2006", "12/05/2023")
		require.NoError(t, err)
		pr := pkg.Product{
			ID:          1,
			Name:        "Tv Samsung",
			Price:       1355.8,
			Quantity:    2,
			CodeValue:   "TV_EXT_SAMSUNG",
			IsPublished: false,
			Expiration:  tim,
		}
		b, err = json.Marshal(pr)
		require.NoError(t, err)
		require.JSONEq(t, string(b), res.Body.String())
		require.Equal(t, 201, res.Code)
	})
	t.Run("should fail due to wrong date", func(t *testing.T) {
		p := pkg.CreateProductRequest{
			Name:        "Tv Samsung",
			Price:       1355.8,
			Quantity:    0,
			CodeValue:   "TV_EXT_SAMSUNG",
			IsPublished: false,
			Expiration:  "12/05/20s23",
		}
		b, err := json.Marshal(p)
		require.NoError(t, err)
		req := httptest.NewRequest("POST", "/product", strings.NewReader(string(b)))
		res := httptest.NewRecorder()
		handler.Create(res, req)
		require.Equal(t, 400, res.Code)
		errReceived, expectedError := GetErrorReceivedAndExpected(t, pkg.WrongFieldValues, res)
		require.Equal(t, errReceived, expectedError)
	})
	t.Run("should fail due to wrong request", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/product", strings.NewReader(`{"other":"Tv Samsung","field":1355.8,"that":2,"does":"TV_EXT_SAMSUNG","not":false,"match":"12/05/2023"}`))
		res := httptest.NewRecorder()
		handler.Create(res, req)
		require.Equal(t, 400, res.Code)
		errorResponse, expectedError := GetErrorReceivedAndExpected(t, pkg.WrongFieldValues, res)
		require.Equal(t, expectedError, errorResponse)
	})
	close(audit)
}

func TestProductGet(t *testing.T) {
	audit := make(chan audit.AuditLog)
	r := chi.NewRouter()
	db := test.NewFakeDb()
	db.FakeMap = map[int]*pkg.Product{
		1: {
			ID:          1,
			Name:        "Tv Samsung",
			Price:       1355.8,
			Quantity:    2,
			CodeValue:   "TV_EXT_SAMSUNG",
			IsPublished: false,
			Expiration:  time.Date(2023, 5, 12, 0, 0, 0, 0, time.UTC),
		},
		2: {
			ID:          2,
			Name:        "Tv LG",
			Price:       15.8,
			Quantity:    1,
			CodeValue:   "TV_EXT_LG",
			IsPublished: false,
			Expiration:  time.Date(2023, 5, 12, 0, 0, 0, 0, time.UTC),
		},
	}
	svc := service.NewProductService(db)
	handler := NewHandler(svc, audit)
	r.Get("/product/{id}", handler.GetByID)
	r.Get("/product", handler.GetAll)

	t.Run("should get a product by the id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/product/1", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		tim, err := time.Parse("02/01/2006", "12/05/2023")
		require.NoError(t, err)
		pr := pkg.Product{
			ID:          1,
			Name:        "Tv Samsung",
			Price:       1355.8,
			Quantity:    2,
			CodeValue:   "TV_EXT_SAMSUNG",
			IsPublished: false,
			Expiration:  tim,
		}
		b, err := json.Marshal(pr)
		require.NoError(t, err)
		require.JSONEq(t, string(b), res.Body.String())
		require.Equal(t, 200, res.Code)
	})
	t.Run("Should retreive the list of all the items", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/product", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		require.Equal(t, 200, res.Code)
	})

	t.Run("should fail due to wrong id format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/product/naa", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		require.Equal(t, 400, res.Code)
		errorResponse, expectedError := GetErrorReceivedAndExpected(t, pkg.BadRequest, res)
		require.Equal(t, expectedError, errorResponse)
	})
	t.Run("should fail due to not found error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/product/7", nil)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		require.Equal(t, 404, res.Code)
		errorResponse, expectedError := GetErrorReceivedAndExpected(t, pkg.NotFound, res)
		require.Equal(t, expectedError, errorResponse)
	})
	close(audit)
}
