package enums

type OrderStatus int32

const (
	OrderStatusPending   OrderStatus = iota
	OrderStatusShipped               = iota
	OrderStatusDelivered             = iota
	OrderStatusCancelled             = iota
)

func (s OrderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "Pending"
	case OrderStatusShipped:
		return "Shipped"
	case OrderStatusDelivered:
		return "Delivered"
	case OrderStatusCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}
