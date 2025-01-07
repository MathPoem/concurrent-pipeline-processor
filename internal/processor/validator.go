package processor

import (
	"concurrent-pipeline-processor/pkg/models"
	"errors"
)

var (
	ErrNilOperations   = errors.New("operations slice is nil")
	ErrInvalidOperator = errors.New("invalid operator")
)

// ValidateTask validates a task and its operations
func ValidateTask(task models.Task) error {
	if task.Operations == nil {
		return ErrNilOperations
	}

	for _, op := range task.Operations {
		if err := validateOperation(op); err != nil {
			return err
		}
	}

	return nil
}

func validateOperation(op models.Operation) error {
	if op.Operator < 0 || op.Operator >= models.OperatorTotalAmount {
		return ErrInvalidOperator
	}

	// For division, check if divisor is zero
	if op.Operator == models.OperatorDivide && op.Value == 0 {
		return errors.New("division by zero")
	}

	return nil
}
