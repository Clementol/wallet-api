package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	User_id      string        `json:"user_id"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID               primitive.ObjectID `bson:"_id"`
	Amount           float64            `json:"amount"`
	From             string             `json:"from"`
	To               string             `json:"to"`
	Transaction_type string             `json:"transaction_type"`
	Created_at       time.Time          `json:"created_at"`
	Transaction_id   string             `json:"transaction_id"`
}
