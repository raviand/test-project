package memory

import (
	"errors"

	"github.com/raviand/test-project/internal/domain"
	"github.com/raviand/test-project/internal/repository"
)

var (
	RecordAlreadyExist = errors.New("conflict, id of record already exist")
	NotFoundError      = errors.New("no record found")
)

type Data struct {
	database map[int]*domain.Product
}

func NewDatabase(filePath string) repository.DataInterface {
	db := make(map[int]*domain.Product)
	return &Data{
		database: db,
	}
}

func (d *Data) CreateProduct(product *domain.Product) (*domain.Product, error) {
	_, exist := d.database[product.Id]
	if exist {
		return nil, RecordAlreadyExist
	}
	d.database[product.Id] = product
	return product, nil
}

func (d *Data) GetProductByCode(code string) (*domain.Product, error) {
	for _, p := range d.database {
		if p.CodeValue == code {
			return p, nil
		}
	}
	return nil, NotFoundError
}

func (d *Data) GetProductById(id int) (*domain.Product, error) {
	p, ok := d.database[id]
	if ok {
		return p, nil
	}
	return nil, NotFoundError
}

func (d *Data) PatchProduct(id int, product *domain.Product) (*domain.Product, error) {
	_, ok := d.database[id]
	if !ok {
		return nil, NotFoundError
	}
	d.database[id] = product
	return d.database[id], nil
}
