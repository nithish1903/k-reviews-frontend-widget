package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SiteReview represents a review for a site.
type SiteReview struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` // Primary key
	AccountID string             `json:"account_id" bson:"account_id"`
	User      User               `json:"user" bson:"user"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Address   Address            `json:"address" bson:"address"`
	Rating    int                `json:"rating" bson:"rating"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	Status    string             `json:"status" bson:"status"`
	VoteID    primitive.ObjectID `json:"vote_id" bson:"vote_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
