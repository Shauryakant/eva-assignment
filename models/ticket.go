package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusClosed     Status = "closed"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusClosed:
		return true
	default:
		return false
	}
}

type Ticket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Status      Status             `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CreateTicketInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTicketStatusInput struct {
	Status Status `json:"status"`
}

func ValidateStatusTransition(currentStatus, newStatus Status) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid ticket status '%s'", newStatus)
	}

	if currentStatus == StatusClosed {
		return fmt.Errorf("closed ticket cannot be updated or reopened")
	}

	if currentStatus == StatusOpen {
		if newStatus != StatusInProgress {
			return fmt.Errorf("invalid transition from 'open' to '%s'; status must progress sequentially (open -> in_progress -> closed)", newStatus)
		}
		return nil
	}

	if currentStatus == StatusInProgress {
		if newStatus != StatusClosed {
			return fmt.Errorf("invalid transition from 'in_progress' to '%s'; status must progress sequentially (open -> in_progress -> closed)", newStatus)
		}
		return nil
	}

	return fmt.Errorf("invalid current status '%s'", currentStatus)
}
