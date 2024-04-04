package handlers

type operation int

const (
	OperationSave operation = iota
	OperationDelete
	OperationUpdate
	OperationGet
)

func (o operation) String() string {
	switch o {
	case OperationSave:
		return "saving"
	case OperationDelete:
		return "deleting"
	case OperationUpdate:
		return "updating"
	case OperationGet:
		return "fetching"
	default:
		return "using"
	}
}
