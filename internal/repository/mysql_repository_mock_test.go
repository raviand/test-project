package repository_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock" // Importa el paquete sqlmock
)

// CreateProductType inserta un nuevo tipo de producto en la base de datos.
// Retorna el ID del tipo de producto existente o recién creado.
func CreateMockProductType(db *sql.DB, name string) (*ProductType, error) {
	// Primero, intenta buscar si el tipo de producto ya existe
	var existingID int
	err := db.QueryRow("SELECT id FROM product_types WHERE name = ?", name).Scan(&existingID)
	if err == nil {
		// Si existe, retorna el tipo de producto existente
		return &ProductType{ID: existingID, Name: name}, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("error al buscar tipo de producto existente: %w", err)
	}

	// Si no existe, inserta uno nuevo
	stmt, err := db.Prepare("INSERT INTO product_types(name) VALUES(?)")
	if err != nil {
		return nil, fmt.Errorf("error al preparar la sentencia INSERT para product_types: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(name)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la sentencia INSERT para product_types: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener el último ID insertado para product_types: %w", err)
	}

	return &ProductType{ID: int(lastID), Name: name}, nil
}

// CreateProduct inserta un nuevo producto en la base de datos.
func CreateMockProduct(db *sql.DB, product *Product) (*Product, error) {
	// Asegúrate de que el ProductType exista o créalo.
	productType, err := CreateMockProductType(db, product.ProductType.Name)
	if err != nil {
		return nil, fmt.Errorf("error al asegurar el tipo de producto: %w", err)
	}
	product.ProductTypeId = productType.ID // Asigna el ID del tipo de producto
	product.ProductType = *productType     // Asigna el objeto ProductType completo

	stmt, err := db.Prepare(`
		INSERT INTO products(name, price, quantity, code_value, is_published, expiration, product_type_id)
		VALUES(?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return nil, fmt.Errorf("error al preparar la sentencia INSERT para products: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		product.Name,
		product.Price,
		product.Quantity,
		product.CodeValue,
		product.IsPublished,
		product.Expiration.Format("2006-01-02"), // Formatear a DATE para MySQL
		product.ProductTypeId,
	)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la sentencia INSERT para products: %w", err)
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error al obtener el último ID insertado para products: %w", err)
	}

	// Recuperar el producto recién creado para obtener el ID y el timestamp generado por la BD.
	// Hacemos un JOIN para obtener también el nombre del tipo de producto.
	row := db.QueryRow(`
		SELECT p.id, p.name, p.price, p.quantity, p.code_value, p.is_published, p.expiration, p.product_type_id, pt.name
		FROM products p
		JOIN product_types pt ON p.product_type_id = pt.id
		WHERE p.id = ?`, lastID)

	var t Product
	var productTypeName string
	if err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Price,
		&t.Quantity,
		&t.CodeValue,
		&t.IsPublished,
		&t.Expiration,
		&t.ProductTypeId,
		&productTypeName,
	); err != nil {
		return nil, fmt.Errorf("error al escanear el producto recién creado: %w", err)
	}
	t.ProductType = ProductType{ID: t.ProductTypeId, Name: productTypeName}

	return &t, nil
}

// GetAllProducts recupera todos los productos de la base de datos, incluyendo su tipo de producto.
func GetAllMockProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query(`
		SELECT p.id, p.name, p.price, p.quantity, p.code_value, p.is_published, p.expiration, p.product_type_id, pt.name
		FROM products p
		JOIN product_types pt ON p.product_type_id = pt.id
		ORDER BY p.id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta SELECT para products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var productTypeName string
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Quantity,
			&p.CodeValue,
			&p.IsPublished,
			&p.Expiration,
			&p.ProductTypeId,
			&productTypeName,
		); err != nil {
			return nil, fmt.Errorf("error al escanear fila de producto: %w", err)
		}
		p.ProductType = ProductType{ID: p.ProductTypeId, Name: productTypeName}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error en las filas de la consulta de productos: %w", err)
	}

	return products, nil
}

// TestProductOperations es un ejemplo de cómo usar go-sqlmock con la entidad Product.
func TestMockProductOperations(t *testing.T) {
	// Crea una conexión mockeada y una instancia de sqlmock
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp)) // Usamos QueryMatcherRegexp
	if err != nil {
		t.Fatalf("error al crear el mock de la base de datos: %v", err)
	}
	defer db.Close()

	t.Run("CreateProduct", func(t *testing.T) {
		// Resetear expectativas para esta subprueba
		mock.ExpectationsWereMet()          // Asegura que las expectativas anteriores se hayan cumplido
		mock.MatchExpectationsInOrder(true) // Opcional: asegura que las expectativas se cumplan en el orden definido

		// Mock para CreateProductType: primero, no encuentra el tipo, luego lo inserta
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM product_types WHERE name = ?")).
			WithArgs("Electronics").
			WillReturnError(sql.ErrNoRows) // Simula que no existe

		mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO product_types(name) VALUES(?)")).
			ExpectExec().
			WithArgs("Electronics").
			WillReturnResult(sqlmock.NewResult(1, 1)) // ID 1 para 'Electronics'

		// Mock para CreateProduct: preparación y ejecución del INSERT
		mock.ExpectPrepare(regexp.QuoteMeta(`
		INSERT INTO products(name, price, quantity, code_value, is_published, expiration, product_type_id)
		VALUES(?, ?, ?, ?, ?, ?, ?)
		`)).
			ExpectExec().
			WithArgs(
				"Laptop",
				1200.50,
				10,
				"LAPTOP001",
				true,
				sqlmock.AnyArg(), // Usamos AnyArg para la fecha, ya que el formato exacto puede variar
				1,                // product_type_id para 'Electronics'
			).
			WillReturnResult(sqlmock.NewResult(101, 1)) // ID 101 para el producto

		// Mock para la recuperación del producto recién creado (SELECT con JOIN)
		rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity", "code_value", "is_published", "expiration", "product_type_id", "pt.name"}).
			AddRow(101, "Laptop", 1200.50, 10, "LAPTOP001", true, time.Now().AddDate(1, 0, 0), 1, "Electronics")

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT p.id, p.name, p.price, p.quantity, p.code_value, p.is_published, p.expiration, p.product_type_id, pt.name
		FROM products p
		JOIN product_types pt ON p.product_type_id = pt.id
		WHERE p.id = ?
		`)).
			WithArgs(101).
			WillReturnRows(rows)

		// Ejecuta la función a testear
		product := &Product{
			Name:        "Laptop",
			Price:       1200.50,
			Quantity:    10,
			CodeValue:   "LAPTOP001",
			IsPublished: true,
			Expiration:  time.Now().AddDate(1, 0, 0),
			ProductType: ProductType{Name: "Electronics"}, // Solo el nombre es suficiente para CreateProductType
		}

		createdProduct, err := CreateMockProduct(db, product)
		if err != nil {
			t.Fatalf("error al crear el producto: %v", err)
		}

		if createdProduct.ID == 0 {
			t.Errorf("ID del producto esperado > 0, obtenido %d", createdProduct.ID)
		}
		if createdProduct.Name != product.Name {
			t.Errorf("Nombre del producto esperado '%s', obtenido '%s'", product.Name, createdProduct.Name)
		}
		if createdProduct.ProductType.Name != "Electronics" {
			t.Errorf("Nombre del tipo de producto esperado 'Electronics', obtenido '%s'", createdProduct.ProductType.Name)
		}

		// Verifica que todas las expectativas se hayan cumplido
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectativas no cumplidas: %s", err)
		}
	})

	t.Run("GetAllProducts", func(t *testing.T) {
		// Resetear expectativas para esta subprueba
		mock.ExpectationsWereMet()
		mock.MatchExpectationsInOrder(true)

		// Mock para CreateProductType (Electronics): ya existe
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM product_types WHERE name = ?")).
			WithArgs("Electronics").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // ID 1 para 'Electronics'

		// Mock para CreateProduct (Go Programming):
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM product_types WHERE name = ?")).
			WithArgs("Books").
			WillReturnError(sql.ErrNoRows) // Simula que no existe

		mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO product_types(name) VALUES(?)")).
			ExpectExec().
			WithArgs("Books").
			WillReturnResult(sqlmock.NewResult(2, 1)) // ID 2 para 'Books'

		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO products(name, price, quantity, code_value, is_published, expiration, product_type_id) VALUES(?, ?, ?, ?, ?, ?, ?)`)).
			ExpectExec().
			WithArgs(
				"Go Programming",
				35.99,
				50,
				"GOBOOK001",
				true,
				sqlmock.AnyArg(),
				2, // product_type_id para 'Books'
			).
			WillReturnResult(sqlmock.NewResult(102, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.name, p.price, p.quantity, p.code_value, p.is_published, p.expiration, p.product_type_id, pt.name FROM products p JOIN product_types pt ON p.product_type_id = pt.id WHERE p.id = ?`)).
			WithArgs(102).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "quantity", "code_value", "is_published", "expiration", "product_type_id", "pt.name"}).
				AddRow(102, "Go Programming", 35.99, 50, "GOBOOK001", true, time.Now().AddDate(2, 0, 0), 2, "Books"))

		// Mock para CreateProduct (Monitor):
		mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO products(name, price, quantity, code_value, is_published, expiration, product_type_id) VALUES(?, ?, ?, ?, ?, ?, ?)`)).
			ExpectExec().
			WithArgs(
				"Monitor",
				299.99,
				20,
				"MONITOR001",
				false,
				sqlmock.AnyArg(),
				1, // product_type_id para 'Electronics'
			).
			WillReturnResult(sqlmock.NewResult(103, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT p.id, p.name, p.price, p.quantity, p.code_value, p.is_published, p.expiration, p.product_type_id, pt.name FROM products p JOIN product_types pt ON p.product_type_id = pt.id WHERE p.id = ?`)).
			WithArgs(103).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "price", "quantity", "code_value", "is_published", "expiration", "product_type_id", "pt.name"}).
				AddRow(103, "Monitor", 299.99, 20, "MONITOR001", false, time.Now().AddDate(0, 6, 0), 1, "Books"))

		// Ejecuta la creación de productos para poblar el mock
		_, err = CreateMockProduct(db, &Product{
			Name:        "Go Programming",
			Price:       35.99,
			Quantity:    50,
			CodeValue:   "GOBOOK001",
			IsPublished: true,
			Expiration:  time.Now().AddDate(2, 0, 0),
			ProductType: ProductType{Name: "Books"},
		})
		if err != nil {
			t.Fatalf("error al crear el segundo producto: %v", err)
		}

		_, err = CreateMockProduct(db, &Product{
			Name:        "Monitor",
			Price:       299.99,
			Quantity:    20,
			CodeValue:   "MONITOR001",
			IsPublished: false,
			Expiration:  time.Now().AddDate(0, 6, 0),
			ProductType: ProductType{Name: "Electronics"},
		})
		if err != nil {
			t.Fatalf("error al crear el tercer producto: %v", err)
		}

		// Mock para GetAllProducts
		rows := sqlmock.NewRows([]string{"id", "name", "price", "quantity", "code_value", "is_published", "expiration", "product_type_id", "pt.name"}).
			AddRow(101, "Laptop", 1200.50, 10, "LAPTOP001", true, time.Now().AddDate(1, 0, 0), 1, "Electronics").
			AddRow(102, "Go Programming", 35.99, 50, "GOBOOK001", true, time.Now().AddDate(2, 0, 0), 2, "Books").
			AddRow(103, "Monitor", 299.99, 20, "MONITOR001", false, time.Now().AddDate(0, 6, 0), 1, "Electronics")

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT p.id, p.name, p.price, p.quantity, p.code_value, p.is_published, p.expiration, p.product_type_id, pt.name
		FROM products p
		JOIN product_types pt ON p.product_type_id = pt.id
		ORDER BY p.id ASC
		`)).
			WillReturnRows(rows)

		allProducts, err := GetAllMockProducts(db)
		if err != nil {
			t.Fatalf("error al obtener todos los productos: %v", err)
		}

		expectedCount := 3
		if len(allProducts) != expectedCount {
			t.Errorf("Cantidad de productos esperada %d, obtenida %d", expectedCount, len(allProducts))
		}

		// Verifica que todas las expectativas se hayan cumplido
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectativas no cumplidas: %s", err)
		}
	})

	// Asegúrate de que no haya expectativas pendientes al final de todos los tests
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas pendientes al final de los tests: %s", err)
	}
}
