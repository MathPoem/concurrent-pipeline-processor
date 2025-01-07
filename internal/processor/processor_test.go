package processor

import (
	"testing"

	"concurrent-pipeline-processor/pkg/models"
)

func TestProcessTask(t *testing.T) {
	tests := []struct {
		name    string
		task    models.Task
		want    int
		wantErr bool
		errMsg  string
	}{
		{
			name: "basic addition",
			task: models.Task{
				Value: 5,
				Operations: []models.Operation{
					{Operator: models.OperatorPlus, Value: 3},
				},
			},
			want: 8,
		},
		{
			name: "multiple operations",
			task: models.Task{
				Value: 10,
				Operations: []models.Operation{
					{Operator: models.OperatorPlus, Value: 5},
					{Operator: models.OperatorMultiply, Value: 2},
					{Operator: models.OperatorMinus, Value: 3},
				},
			},
			want: 27, // (10 + 5) * 2 - 3
		},
		{
			name: "division by zero",
			task: models.Task{
				Value: 10,
				Operations: []models.Operation{
					{Operator: models.OperatorDivide, Value: 0},
				},
			},
			wantErr: true,
			errMsg:  "division by zero",
		},
		{
			name: "division overflow check",
			task: models.Task{
				Value: -1 << 31, // Minimum int32 value
				Operations: []models.Operation{
					{Operator: models.OperatorDivide, Value: -1},
				},
			},
			wantErr: true,
			errMsg:  "division would cause integer overflow",
		},
		{
			name: "valid division",
			task: models.Task{
				Value: 10,
				Operations: []models.Operation{
					{Operator: models.OperatorDivide, Value: 2},
				},
			},
			want: 5,
		},
		{
			name: "invalid operator",
			task: models.Task{
				Value: 10,
				Operations: []models.Operation{
					{Operator: models.OperatorTotalAmount, Value: 5},
				},
			},
			wantErr: true,
			errMsg:  "invalid operator",
		},
		{
			name: "complex calculation with division",
			task: models.Task{
				Value: 100,
				Operations: []models.Operation{
					{Operator: models.OperatorDivide, Value: 2},   // 50
					{Operator: models.OperatorPlus, Value: 10},    // 60
					{Operator: models.OperatorMultiply, Value: 3}, // 180
					{Operator: models.OperatorMinus, Value: 30},   // 150
				},
			},
			want: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessTask(tt.task)

			if tt.wantErr {
				if result.Error == nil {
					t.Errorf("ProcessTask() expected error %v, got nil", tt.errMsg)
					return
				}
				if result.Error.Error() != tt.errMsg {
					t.Errorf("ProcessTask() error = %v, want %v", result.Error, tt.errMsg)
				}
				return
			}

			if result.Error != nil {
				t.Errorf("ProcessTask() unexpected error: %v", result.Error)
				return
			}

			if result.Result != tt.want {
				t.Errorf("ProcessTask() = %v, want %v", result.Result, tt.want)
			}
		})
	}
}
