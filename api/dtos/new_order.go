package dtos

type NewOrderDto struct {
	ClientName  string
	ClientEmail string
	ClientPhone string
	Note        string
	Items       []OrderItemDto
}
