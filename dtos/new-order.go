package dtos

type NewOrderDto struct {
	ClientName  string
	ClientEmail string
	ClientPhone string
	Note        string
	Items       []OrderItemDto
}

func (order *NewOrderDto) GetTotalAmount() float32 {
	return 0.0
}

func (order *NewOrderDto) GetTaxAmount() float32 {
	return 0.0
}
