package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Task struct {
		ID          interface{}        `json:"_id,omitempty" bson:"_id,omitempty" `
		Title       string             `json:"title" bson:"title" validate:"required"`
		Description string             `json:"description" bson:"description" validate:"required"`
		DueDate     string             `json:"dueDate" bson:"dueDate" validate:"required"`
		Status      string             `json:"status,omitempty" bson:"status" validate:"required"`
		CreatedBy   primitive.ObjectID `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
		AssignedTo  primitive.ObjectID `json:"assignedTo,omitempty" bson:"assignedTo,omitempty"`
	}
)
