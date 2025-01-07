package processor

import (
	"testing"

	"concurrent-pipeline-processor/pkg/models"
)

func TestValidateTask(t *testing.T) {
	tests := []struct {
		name    string
		task    models.Task
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid task",
			task: models.Task{
				Value: 5,
				Operations: []models.Operation{
					{Operator: models.OperatorPlus, Value: 3},
				},
			},
			wantErr: false,
		},
		{
			name: "nil operations",
			task: models.Task{
				Value:      5,
				Operations: nil,
			},
			wantErr: true,
			errMsg:  "operations slice is nil",
		},
		{
			name: "invalid operator",
			task: models.Task{
				Value: 5,
				Operations: []models.Operation{
					{Operator: models.OperatorTotalAmount, Value: 3},
				},
			},
			wantErr: true,
			errMsg:  "invalid operator",
		},
		{
			name: "division by zero",
			task: models.Task{
				Value: 5,
				Operations: []models.Operation{
					{Operator: models.OperatorDivide, Value: 0},
				},
			},
			wantErr: true,
			errMsg:  "division by zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTask(tt.task)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateTask() expected error %v, got nil", tt.errMsg)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("ValidateTask() error = %v, want %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateTask() unexpected error: %v", err)
			}
		})
	}
}
