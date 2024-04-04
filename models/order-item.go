package models

type OrderItem struct {
	ID        int64
	UnitPrice float32 `binding:"required"`
	Tax       float32 `binding:"required"`
	Quantity  int64   `binding:"required"`
	OrderID   int64   `binding:"required"`
	ProductID int64   `binding:"required"`
}
