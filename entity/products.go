package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProductInfo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	AccountID string             `json:"account_id" bson:"account_id"`
	ProductID int64              `json:"product_id" bson:"product_id"`
	Name      string             `json:"name" bson:"name"`
	ImageURL  string             `json:"image_url" bson:"image_url"`
	PageURL   string             `json:"page_url" bson:"page_url"`
	Price     string             `json:"price,omitempty" bson:"price,omitempty"`
	Status    string             `json:"status" bson:"status"`
}
