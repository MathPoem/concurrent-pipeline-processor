package models

type Operator int

const (
	OperatorPlus Operator = iota
	OperatorMinus
	OperatorDivide
	OperatorMultiply
	OperatorTotalAmount
)

type Operation struct {
	Operator Operator
	Value    int
}

type Task struct {
	Value      int
	Operations []Operation
}

type Result struct {
	Result int
	Error  error
}
