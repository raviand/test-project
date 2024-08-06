package repository

import (
	"fmt"
	"os"

	"github.com/raviand/test-project/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type gormDb struct {
	db *gorm.DB
}

func NewGormDb() Database {
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/movies_db?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"))
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = database.AutoMigrate(&pkg.Product{}, &pkg.ProductType{})
	if err != nil {
		panic("failed to migrate database")
	}
	return &gormDb{db: database}
}

func (g *gormDb) ListAll() ([]*pkg.Product, error) {
	var products []*pkg.Product
	g.db.Find(&products)
	return products, nil
}

func (g *gormDb) GetByID(id int) (*pkg.Product, error) {
	var product pkg.Product
	g.db.First(&product, id)
	return &product, nil
}

func (g *gormDb) Create(product *pkg.Product) error {
	var productType pkg.ProductType
	g.db.First(&productType, product.ProductTypeId)
	if productType.ID == 0 {
		return pkg.GetError(pkg.BadRequest)
	}
	product.ProductType = productType
	g.db.Create(product)
	return nil
}

func (g *gormDb) Update(product *pkg.Product) error {
	g.db.Save(product)
	return nil
}

func (g *gormDb) Delete(id int) error {
	g.db.Delete(&pkg.Product{}, id)
	return nil
}
