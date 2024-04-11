package dtos

type NewOrderDto struct {
	ClientName  string
	ClientEmail string
	ClientPhone string
	Note        string
	Items       []OrderItemDto
}

type OrderItemDto struct {
	Quantity  int32
	ProductID int32
}
