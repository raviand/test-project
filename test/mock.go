package test

import "github.com/raviand/test-project/pkg"

type fakeDb struct {
	FakeMap map[int]*pkg.Product
	Fail    bool
}

func NewFakeDb() *fakeDb {
	return &fakeDb{
		FakeMap: make(map[int]*pkg.Product),
	}
}

func (d *fakeDb) ListAll() ([]*pkg.Product, error) {
	products := make([]*pkg.Product, 0, len(d.FakeMap))
	for _, product := range d.FakeMap {
		products = append(products, product)
	}
	return products, nil
}

func (d *fakeDb) GetByID(id int) (*pkg.Product, error) {
	product, ok := d.FakeMap[id]
	if !ok {
		return nil, pkg.GetError(pkg.NotFound)
	}
	return product, nil
}

func (d *fakeDb) Create(product *pkg.Product) error {
	if d.Fail {
		return pkg.GetError(pkg.InternalError)
	}
	d.FakeMap[product.ID] = product
	return nil
}

func (d *fakeDb) Update(product *pkg.Product) error {
	if _, ok := d.FakeMap[product.ID]; !ok {
		return pkg.GetError(pkg.NotFound)
	}
	d.FakeMap[product.ID] = product
	return nil
}

func (d *fakeDb) Delete(id int) error {
	if _, ok := d.FakeMap[id]; !ok {
		return pkg.GetError(pkg.NotFound)
	}
	delete(d.FakeMap, id)
	return nil
}
