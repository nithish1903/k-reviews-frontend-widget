package entity

import (
	"time"
)

type SaveProductReviewRequest struct {
	ProductReview ProductReview `json:"product_review"`
	OtherField1   string        `json:"other_field1"`
	OtherField2   int           `json:"other_field2"`
	// Add other fields as needed
}

// ProductReview represents a review for a product.
type ProductReview struct {
	// Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	ProductID   int64        `json:"product_id" bson:"product_id"`
	AccountID   string       `json:"account_id" bson:"account_id"`
	User        User         `json:"user" bson:"user"`
	Rating      int          `json:"rating" bson:"rating"`
	Title       string       `json:"title" bson:"title"`
	Content     string       `json:"content" bson:"content"`
	MediaUpload MediaUploads `json:"media_uploads" bson:"media_uploads"`
	Status      string       `json:"status" bson:"status"`
	VoteID      string       `json:"vote_id" bson:"vote_id"`
	CreatedAt   time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" bson:"updated_at"`
}

// User represents a user information.
type User struct {
	Name    string  `json:"name" bson:"name"`
	Email   string  `json:"email" bson:"email"`
	Address Address `json:"address" bson:"address"`
}

// Address represents user's address details.
type Address struct {
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Country string `json:"country" bson:"country"`
	Pincode string `json:"pincode" bson:"pincode"`
}

// MediaUpload represents uploaded media files.
type MediaUploads struct {
	Image []Image `json:"image" bson:"image"`
	Video []Video `json:"video" bson:"video"`
}

// Image represents an image file.
type Image struct {
	Name      string `json:"name" bson:"name"`
	Status    string `json:"status" bson:"status"`
	SortOrder int    `json:"sort_order" bson:"sort_order"`
}

// Video represents a video file.
type Video struct {
	Name      string `json:"name" bson:"name"`
	Status    string `json:"status" bson:"status"`
	SortOrder int    `json:"sort_order" bson:"sort_order"`
}

type ReviewResponse struct {
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Status     bool        `json:"status,omitempty"`
	StatusCode int         `json:"statuscode,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Timestamp  time.Time   `json:"timestamp,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type PaginationResponse struct {
	Message    string                   `json:"message,omitempty"`
	Data       []map[string]interface{} `json:"data,omitempty"`
	Status     bool                     `json:"status,omitempty"`
	StatusCode int                      `json:"statuscode,omitempty"`
	Error      interface{}              `json:"error,omitempty"`
	Timestamp  time.Time                `json:"timestamp,omitempty"`
	Pagination *Pagination              `json:"pagination,omitempty"`
}
type Pagination struct {
	Total   int64 `json:"total,omitempty"`
	PerPage int64 `json:"per_page,omitempty"`
	Page    int64 `json:"page,omitempty"`
}
