package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductQuestionAnswer represents a question-answer pair for a product.
type ProductQuestionAnswer struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` // Primary key
	AccountID string             `json:"account_id" bson:"account_id"`
	ProductID int64              `json:"product_id" bson:"product_id"`
	User      string             `json:"user" bson:"user"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Address   Address            `json:"address" bson:"address"`
	Question  string             `json:"question" bson:"question"`
	Answer    []Answer           `json:"answer" bson:"answer"`
	VoteID    primitive.ObjectID `json:"vote_id" bson:"vote_id"`
	Status    string             `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type Answer struct {
	ID        string    `json:"ans_id" bson:"ans_id"`
	Answer    string    `json:"answer" bson:"answer"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Upvotes   int64     `json:"up_votes,omitempty" bson:"up_votes,omitempty"`
	Downvotes int64     `json:"down_votes,omitempty" bson:"down_votes,omitempty"` // Primary key
}

// Votes represents votes on various modules.
// type Votes struct {
// 	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` // Primary key
// 	ModuleID primitive.ObjectID `json:"module_id" bson:"module_id"`
// 	Module   string             `json:"module" bson:"module"`
// 	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
// 	Status   string             `json:"status" bson:"status"`
// }
