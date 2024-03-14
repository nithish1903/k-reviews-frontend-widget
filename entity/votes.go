package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Votes represents a review for a site.
type Votes struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` // Primary key
	AccountID string             `json:"account_id" bson:"account_id"`
	ModuleId  string             `json:"module_id" bson:"module_id"`
	Mudule    string             `json:"module" bson:"module"`
	UserId    string             `json:"user_id" bson:"user_id"`
	Status    int64              `json:"status" bson:"status"`
}
