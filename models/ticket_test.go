package models

import (
	"testing"
)

func TestValidateStatusTransition(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus Status
		newStatus     Status
		expectError   bool
	}{
		{
			name:          "Valid transition open to in_progress",
			currentStatus: StatusOpen,
			newStatus:     StatusInProgress,
			expectError:   false,
		},
		{
			name:          "Valid transition in_progress to closed",
			currentStatus: StatusInProgress,
			newStatus:     StatusClosed,
			expectError:   false,
		},
		{
			name:          "Invalid transition open to closed (jumping steps)",
			currentStatus: StatusOpen,
			newStatus:     StatusClosed,
			expectError:   true,
		},
		{
			name:          "Invalid transition closed to open (reopening)",
			currentStatus: StatusClosed,
			newStatus:     StatusOpen,
			expectError:   true,
		},
		{
			name:          "Invalid transition closed to in_progress",
			currentStatus: StatusClosed,
			newStatus:     StatusInProgress,
			expectError:   true,
		},
		{
			name:          "Invalid transition open to open",
			currentStatus: StatusOpen,
			newStatus:     StatusOpen,
			expectError:   true,
		},
		{
			name:          "Invalid transition in_progress to open (backward)",
			currentStatus: StatusInProgress,
			newStatus:     StatusOpen,
			expectError:   true,
		},
		{
			name:          "Invalid new status string",
			currentStatus: StatusOpen,
			newStatus:     Status("unknown"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStatusTransition(tt.currentStatus, tt.newStatus)
			if tt.expectError && err == nil {
				t.Errorf("expected error for transition %s -> %s, got nil", tt.currentStatus, tt.newStatus)
			}
			if !tt.expectError && err != nil {
				t.Errorf("expected no error for transition %s -> %s, got: %v", tt.currentStatus, tt.newStatus, err)
			}
		})
	}
}
