package repository_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/DATA-DOG/go-txdb"    // Importa el driver txdb
	_ "github.com/go-sql-driver/mysql" // Importa el driver de MySQL para usar con txdb
)

// ProductType representa la estructura de un tipo de producto.
type ProductType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Product representa la estructura de un producto.
type Product struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Price         float64     `json:"price"`
	Quantity      int         `json:"quantity"`
	CodeValue     string      `json:"code_value"`
	IsPublished   bool        `json:"is_published"`
	Expiration    time.Time   `json:"expiration"`
	ProductTypeId int         `json:"product_type_id"`
	ProductType   ProductType `json:"product_type"` // Campo para la relación, no mapeado directamente a la columna
}

// init se ejecuta antes que cualquier función de test.
// Registra el driver "txdb" para usarlo con "mysql".
func init() {
	// Registra el driver txdb con un nombre de driver subyacente y una cadena de conexión DSN.
	// El DSN aquí es para una base de datos MySQL.
	// ¡ADVERTENCIA!: Reemplaza "user:password@tcp(127.0.0.1:3306)/testdb" con tus credenciales y DSN de MySQL reales.
	// Asegúrate de que la base de datos 'testdb' exista y que el usuario tenga permisos.
	txdb.Register("txdb", "mysql", "username:password@tcp(127.0.0.1:3306)/test?parseTime=true")
}

// setupDB inicializa la base de datos y crea las tablas 'product_types' y 'products'.
// Se asegura de que la base de datos esté lista para las operaciones de test.
func setupDB(db *sql.DB) error {
	// Crea la tabla 'product_types' si no existe.
	createProductTypesTableSQL := `
	CREATE TABLE IF NOT EXISTS product_types (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE
	);`
	_, err := db.Exec(createProductTypesTableSQL)
	if err != nil {
		return fmt.Errorf("error al crear la tabla 'product_types': %w", err)
	}

	// Crea la tabla 'products' si no existe, usando sintaxis MySQL.
	createProductsTableSQL := `
	CREATE TABLE IF NOT EXISTS products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(20) NOT NULL,
		price DOUBLE(10,2) NOT NULL,
		quantity INT,
		code_value VARCHAR(20),
		is_published BOOLEAN,
		expiration DATE,
		product_type_id INT NOT NULL,
		FOREIGN KEY (product_type_id) REFERENCES product_types(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(createProductsTableSQL)
	if err != nil {
		return fmt.Errorf("error al crear la tabla 'products': %w", err)
	}
	return nil
}

// CreateProductType inserta un nuevo tipo de producto en la base de datos.
// Retorna el ID del tipo de producto existente o recién creado.
func CreateProductType(db *sql.DB, name string) (*ProductType, error) {
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
func CreateProduct(db *sql.DB, product *Product) (*Product, error) {
	// Asegúrate de que el ProductType exista o créalo.
	productType, err := CreateProductType(db, product.ProductType.Name)
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
func GetAllProducts(db *sql.DB) ([]Product, error) {
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

// TestProductOperations es un ejemplo de cómo usar go-txdb con la entidad Product.
func TestProductOperations(t *testing.T) {
	// Abre una conexión a la base de datos usando el driver "txdb".
	db, err := sql.Open("txdb", "test_db_connection_products") // DSN arbitrario para este test
	if err != nil {
		t.Fatalf("no se pudo abrir la conexión a la base de datos: %v", err)
	}
	defer db.Close() // Esto hará un ROLLBACK de la transacción al final del test.

	// Configura la base de datos (crea las tablas) para este test.
	if err := setupDB(db); err != nil {
		t.Fatalf("error al configurar la base de datos para el test: %v", err)
	}

	// --- Prueba de la operación CreateProduct ---
	t.Run("CreateProduct", func(t *testing.T) {
		// Crea un tipo de producto
		electronicsType, err := CreateProductType(db, "Electronics")
		if err != nil {
			t.Fatalf("error al crear tipo de producto 'Electronics': %v", err)
		}

		product := &Product{
			Name:        "Laptop",
			Price:       1200.50,
			Quantity:    10,
			CodeValue:   "LAPTOP001",
			IsPublished: true,
			Expiration:  time.Now().AddDate(1, 0, 0), // Expira en 1 año
			ProductType: *electronicsType,
		}

		createdProduct, err := CreateProduct(db, product)
		if err != nil {
			t.Fatalf("error al crear el producto: %v", err)
		}

		if createdProduct.ID == 0 {
			t.Errorf("ID del producto esperado > 0, obtenido %d", createdProduct.ID)
		}
		if createdProduct.Name != product.Name {
			t.Errorf("Nombre del producto esperado '%s', obtenido '%s'", product.Name, createdProduct.Name)
		}
		if createdProduct.Price != product.Price {
			t.Errorf("Precio del producto esperado %.2f, obtenido %.2f", product.Price, createdProduct.Price)
		}
		if createdProduct.ProductType.ID != electronicsType.ID {
			t.Errorf("ID del tipo de producto esperado %d, obtenido %d", electronicsType.ID, createdProduct.ProductType.ID)
		}
		if createdProduct.ProductType.Name != electronicsType.Name {
			t.Errorf("Nombre del tipo de producto esperado '%s', obtenido '%s'", electronicsType.Name, createdProduct.ProductType.Name)
		}

		t.Logf("Producto creado exitosamente: %+v", createdProduct)
	})

	// --- Prueba de la operación GetAllProducts ---
	t.Run("GetAllProducts", func(t *testing.T) {
		// Crea algunos productos adicionales para asegurar que GetAll funcione correctamente
		// Reutiliza o crea nuevos tipos de producto si es necesario
		booksType, err := CreateProductType(db, "Books")
		if err != nil {
			t.Fatalf("error al crear tipo de producto 'Books': %v", err)
		}

		_, err = CreateProduct(db, &Product{
			Name:        "Go Programming",
			Price:       35.99,
			Quantity:    50,
			CodeValue:   "GOBOOK001",
			IsPublished: true,
			Expiration:  time.Now().AddDate(2, 0, 0),
			ProductType: *booksType,
		})
		if err != nil {
			t.Fatalf("error al crear el segundo producto: %v", err)
		}

		_, err = CreateProduct(db, &Product{
			Name:        "Monitor",
			Price:       299.99,
			Quantity:    20,
			CodeValue:   "MONITOR001",
			IsPublished: false,
			Expiration:  time.Now().AddDate(0, 6, 0),      // Expira en 6 meses
			ProductType: ProductType{Name: "Electronics"}, // Usará el tipo "Electronics" existente
		})
		if err != nil {
			t.Fatalf("error al crear el tercer producto: %v", err)
		}

		allProducts, err := GetAllProducts(db)
		if err != nil {
			t.Fatalf("error al obtener todos los productos: %v", err)
		}

		expectedCount := 3 // Incluyendo el producto creado en la subprueba "CreateProduct"
		if len(allProducts) != expectedCount {
			t.Errorf("Cantidad de productos esperada %d, obtenida %d", expectedCount, len(allProducts))
		}

		// Verifica que los nombres de los productos y sus tipos sean los esperados
		expectedProductNames := map[string]string{
			"Laptop":         "Electronics",
			"Go Programming": "Books",
			"Monitor":        "Electronics",
		}
		for _, product := range allProducts {
			if expectedType, ok := expectedProductNames[product.Name]; ok {
				if product.ProductType.Name != expectedType {
					t.Errorf("Tipo de producto inesperado para '%s': esperado '%s', obtenido '%s'", product.Name, expectedType, product.ProductType.Name)
				}
			} else {
				t.Errorf("Nombre de producto inesperado encontrado: %s", product.Name)
			}
		}

		t.Logf("Productos recuperados exitosamente: %+v", allProducts)
	})

	// Después de este test, la transacción de la base de datos se revertirá (ROLLBACK),
	// asegurando que la base de datos esté en su estado original para el siguiente test.
}
