package processor

import (
	"concurrent-pipeline-processor/pkg/models"
	"errors"
)

var (
	ErrDivisionByZero = errors.New("division by zero")
)

// ProcessTask processes a single task by applying all operations
func ProcessTask(task models.Task) models.Result {
	result := task.Value

	for _, op := range task.Operations {
		var err error
		result, err = applyOperation(result, op)
		if err != nil {
			return models.Result{Error: err}
		}
	}

	return models.Result{Result: result}
}

func applyOperation(value int, op models.Operation) (int, error) {
	switch op.Operator {
	case models.OperatorPlus:
		return value + op.Value, nil
	case models.OperatorMinus:
		return value - op.Value, nil
	case models.OperatorMultiply:
		return value * op.Value, nil
	case models.OperatorDivide:
		if op.Value == 0 {
			return 0, ErrDivisionByZero
		}
		// Check if division would cause integer overflow
		if value == -1<<31 && op.Value == -1 {
			return 0, errors.New("division would cause integer overflow")
		}
		return value / op.Value, nil
	default:
		return 0, ErrInvalidOperator
	}
}
