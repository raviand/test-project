package repository

import "github.com/raviand/test-project/pkg"

type Database interface {
	ListAll() ([]*pkg.Product, error)
	GetByID(id int) (*pkg.Product, error)
	Create(product *pkg.Product) error
	Update(product *pkg.Product) error
	Delete(id int) error
}

type db struct {
	products map[int]*pkg.Product
	lastId   int
}

func NewDatabase() Database {
	return &db{
		products: make(map[int]*pkg.Product),
		lastId:   0,
	}
}

func (d *db) ListAll() ([]*pkg.Product, error) {
	products := make([]*pkg.Product, 0, len(d.products))
	for _, product := range d.products {
		products = append(products, product)
	}
	return products, nil
}

func (d *db) GetByID(id int) (*pkg.Product, error) {
	product, ok := d.products[id]
	if !ok {
		return nil, pkg.GetError(pkg.NotFound)
	}
	return product, nil
}

func (d *db) Create(product *pkg.Product) error {
	d.lastId++
	product.ID = d.lastId
	d.products[d.lastId] = product
	return nil
}

func (d *db) Update(product *pkg.Product) error {
	if _, ok := d.products[product.ID]; !ok {
		return pkg.GetError(pkg.NotFound)
	}
	d.products[product.ID] = product
	return nil
}

func (d *db) Delete(id int) error {
	if _, ok := d.products[id]; !ok {
		return pkg.GetError(pkg.NotFound)
	}
	delete(d.products, id)
	return nil
}
