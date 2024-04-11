// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package datastore

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type NullOrderStatus struct {
	OrderStatus OrderStatus
	Valid       bool // Valid is true if OrderStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullOrderStatus) Scan(value interface{}) error {
	if value == nil {
		ns.OrderStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.OrderStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullOrderStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.OrderStatus), nil
}

type Order struct {
	ID                 int32
	Status             OrderStatus
	TotalAmount        pgtype.Numeric
	Tax                pgtype.Numeric
	UserID             int32
	Created            pgtype.Timestamp
	Shipped            pgtype.Timestamp
	Cancelled          pgtype.Timestamp
	Delivered          pgtype.Timestamp
	ClientName         pgtype.Text
	ClientEmail        pgtype.Text
	ClientPhone        pgtype.Text
	Note               pgtype.Text
	CancellationReason pgtype.Text
}

type OrderItem struct {
	ID        int32
	UnitPrice pgtype.Numeric
	Quantity  int32
	OrderID   int32
	ProductID int32
}

type Product struct {
	ID          int32
	Name        string
	Description string
	Sku         string
	Price       pgtype.Numeric
	Stock       int32
}

type User struct {
	ID       int32
	Name     string
	Code     string
	Email    string
	Password string
}
