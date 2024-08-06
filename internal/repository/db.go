package repository

import (
	"database/sql"
	"errors"

	"github.com/raviand/test-project/pkg"
)

type Database interface {
	ListAll() ([]*pkg.Product, error)
	GetByID(id int) (*pkg.Product, error)
	Create(product *pkg.Product) error
	Update(product *pkg.Product) error
	Delete(id int) error
}

type db struct {
	filePath string
	products map[int]*pkg.Product
	lastId   int
	db       *sql.DB
}

func NewDatabase(mysql *sql.DB) Database {
	return &db{
		products: make(map[int]*pkg.Product),
		lastId:   0,
		db:       mysql,
	}
}

func (d *db) ListAll() ([]*pkg.Product, error) {
	// query
	rows, err := d.db.Query(
		"select p.id, p.name, p.quantity, p.price, p.code_value,  p.is_published from product p",
	)
	if err != nil {
		return nil, pkg.GetError(pkg.InternalError)
	}
	defer rows.Close()
	products := make([]*pkg.Product, 0)
	for rows.Next() {
		var product pkg.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price, &product.CodeValue, &product.IsPublished); err != nil {
			return nil, pkg.GetError(pkg.InternalError)
		}
		products = append(products, &product)
	}
	if nil != rows.Err() {
		return nil, pkg.GetError(pkg.InternalError)
	}
	return products, nil
}

func (d *db) GetByID(id int) (*pkg.Product, error) {
	// execute query
	row := d.db.QueryRow("select p.id, p.name, p.quantity, p.price, p.code_value,  p.is_published from product p where m.id = ?", id)
	if err := row.Err(); err != nil {
		return nil, pkg.GetError(pkg.NotFound)
	}

	// scan result
	var product *pkg.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price, &product.CodeValue, &product.IsPublished); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.GetError(pkg.NotFound)
		}
		return nil, err
	}

	return product, nil
}

func (d *db) Create(product *pkg.Product) error {
	// execute query
	result, err := d.db.Exec("INSERT INTO product (name, price, quantity, code_value, is_published, expiration) VALUES (?, ?, ?, ?, ?, ?)", product.Name, product.Price, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration)
	if err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	// get last inserted id
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	// set user id
	product.ID = int(lastInsertId)
	return nil
}

func (d *db) Update(product *pkg.Product) error {
	// execute query
	_, err := d.db.Exec("update product set name = ?, price = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ? where id = ?", product.Name, product.Price, product.Quantity, product.CodeValue, product.IsPublished, product.Expiration, product.ID)
	if err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	return nil
}

func (d *db) Delete(id int) error {
	// execute query
	_, err := d.db.Exec("delete from product where id = ?", id)
	if err != nil {
		return pkg.GetError(pkg.InternalError)
	}
	return nil
}
