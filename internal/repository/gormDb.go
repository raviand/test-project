package repository

import (
	"errors"
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
	if err := g.db.Preload("ProductType").Find(&products).Error; err != nil {
		return nil, pkg.GetError(pkg.InternalError)
	}
	return products, nil
}

func (g *gormDb) GetByID(id int) (*pkg.Product, error) {
	var product pkg.Product
	if err := g.db.Preload("ProductType").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, pkg.GetError(pkg.NotFound)
		}
		return nil, pkg.GetError(pkg.InternalError)
	}
	return &product, nil
}

func (g *gormDb) Create(product *pkg.Product) error {
	var productType pkg.ProductType
	if err := g.db.First(&productType, product.ProductTypeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkg.GetError(pkg.BadRequest)
		}
		return pkg.GetError(pkg.InternalError)
	}

	product.ProductType = productType
	if err := g.db.Create(product).Error; err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	return nil
}

func (g *gormDb) Update(product *pkg.Product) error {
	if err := g.db.Save(product).Error; err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	return nil
}

func (g *gormDb) Delete(id int) error {
	if err := g.db.Delete(&pkg.Product{}, id).Error; err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	return nil
}
