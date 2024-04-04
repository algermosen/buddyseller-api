package models

import (
	"time"

	"example/buddyseller-api/database"
	"example/buddyseller-api/db"
	"example/buddyseller-api/dtos"
	"example/buddyseller-api/enums"
)

type Order struct {
	ID                 int64
	Status             enums.OrderStatus `binding:"required"`
	TotalAmount        float32           `binding:"required"`
	Tax                float32           `binding:"required"`
	CreatedDate        time.Time         `binding:"required"`
	ShippedDate        *time.Time
	CancelledDate      *time.Time
	DeliveredDate      *time.Time
	ClientName         string
	ClientEmail        string
	ClientPhone        string
	Note               string
	CancellationReason string
}

func GetAllOrders() ([]Order, error) {
	query := "SELECT * FROM orders"
	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order
		err := scanOrder(rows, &order)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func GetOrderById(id int64) (*Order, error) {
	query := "SELECT * FROM orders WHERE id = $1"
	row := database.DB.QueryRow(query, id)

	var order Order
	err := scanOrder(row, &order)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func PlaceOrder(newOrder dtos.NewOrderDto) error {
	saveOrderQuery := `
	INSERT INTO orders(status, total_amount, tax, client_name, client_email, client_phone, note)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
	`

	var orderPk int64

	err := database.DB.QueryRow(saveOrderQuery, enums.OrderStatusPending, newOrder.GetTotalAmount(), newOrder.GetTaxAmount(), &newOrder.ClientName, &newOrder.ClientEmail, &newOrder.ClientPhone, &newOrder.Note).Scan(&orderPk)

	if err != nil {
		return err
	}

	saveOrderItemQuery := `
	INSERT INTO order_items(unit_price, tax, quantity, order_id, product_id)
	VALUES ($1, $2, $3, $4, $5)
	`

	stmt, err := database.DB.Prepare(saveOrderItemQuery)

	if err != nil {
		return err
	}

	defer stmt.Close()
	getProductQuery := `
	SELECT price FROM products WHERE id = $1
	`

	_, err = database.DB.Query(getProductQuery)

	if err != nil {
		return err
	}

	const taxValue = 0.18

	// for _, value := range newOrder.Items {
	// 	stmt.Exec(unit_price, tax, value.Quantity, orderPk, value.ProductID)
	// }

	// TODO: Get all product prices at once and make the proper calculations

	return nil
}

func CancelOrder(id int64) error {
	query := `
	UPDATE orders
	SET
		status = $2,
	WHERE id = $1
	`

	stmt, err := database.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(&id, enums.OrderStatusCancelled)

	if err != nil {
		return err
	}

	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanOrder(row rowScanner, order *Order) error {
	return row.Scan(&order.ID,
		&order.Status,
		&order.TotalAmount,
		&order.Tax,
		&order.CreatedDate,
		&order.CreatedDate,
		&order.ShippedDate,
		&order.CancelledDate,
		&order.DeliveredDate,
		&order.ClientName,
		&order.ClientEmail,
		&order.ClientPhone,
		&order.Note,
		&order.CancellationReason,
	)
}
