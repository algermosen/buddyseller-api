package models

import (
	"example/buddyseller-api/database"
	"fmt"
)

type Product struct {
	ID          int64
	Name        string  `binding:"required"`
	Description string  `binding:"required"`
	Sku         string  `binding:"required"`
	Price       float32 `binding:"required"`
	Stock       int64   `binding:"required"`
}

func (product *Product) Save() error {
	query := `
	INSERT INTO products(name, description, sku, price, stock)
	VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	fmt.Printf("%+v\n", product)
	var pk int64

	err := database.DB.QueryRow(query, &product.Name, &product.Description, &product.Sku, &product.Price, &product.Stock).Scan(&pk)

	if err != nil {
		return err
	}

	product.ID = pk
	return nil

}

func GetAllProducts() ([]Product, error) {
	query := "SELECT * FROM products"
	rows, err := database.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Sku, &product.Price, &product.Stock)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func GetProductById(id int64) (*Product, error) {
	query := "SELECT * FROM products WHERE id = $1"
	row := database.DB.QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Sku, &product.Price, &product.Stock)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func GetBySku(sku string) (*Product, error) {
	query := "SELECT * FROM products WHERE sku = $1"
	row := database.DB.QueryRow(query, sku)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Sku, &product.Price, &product.Stock)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (product *Product) Update() error {
	query := `
	UPDATE products
	SET 
		name = $2,
		description = $3,
		sku = $4,
		price = $5,
		stock = $6
	WHERE id = $1
	`

	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(&product.ID, &product.Name, &product.Description, &product.Sku, &product.Price, &product.Stock)

	if err != nil {
		return err
	}

	return nil
}

func DeleteProduct(id int64) error {
	query := `
	DELETE FROM products
	WHERE id = $1
	`

	_, err := database.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
