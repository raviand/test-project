package pkg

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	ProductType int     `json:"product_type"`
}

type Product struct {
	ID            int         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string      `json:"name" gorm:"type:varchar(20);not null"`
	Price         float64     `json:"price" gorm:"type:double(10,2);not null"`
	Quantity      int         `json:"quantity"`
	CodeValue     string      `json:"code_value" gorm:"type:varchar(20)"`
	IsPublished   bool        `json:"is_published"`
	Expiration    time.Time   `json:"expiration" gorm:"type:date"`
	ProductTypeId int         `json:"product_type_id" gorm:"not null"`
	ProductType   ProductType `json:"product_type" gorm:"foreignKey:ProductTypeID"`
}

type ProductType struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"type:varchar(20);not null"`
}

// Implementación de la interfaz Valuer para ProductType
func (pt ProductType) Value() (driver.Value, error) {
	return json.Marshal(pt)
}

// Implementación de la interfaz Scanner para ProductType
func (pt *ProductType) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan ProductType: %v", value)
	}
	return json.Unmarshal(bytes, pt)
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

type User struct {
	Id         string `json:"id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	IP         string `json:"ip"`
	MacAddress string `json:"macAddress"`
	Website    string `json:"website"`
	Image      string `json:"image"`
}
