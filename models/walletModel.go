package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID         primitive.ObjectID `bson:"_id"`
	User_id    string             `json:"user_id"`
	Balance    float64            `json:"balance"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Wallet_id  string             `json:"wallet_id"`
	Active     *bool              `json:"active"`
}

type FundWallet struct {
	// User_id string  `json:"user_id" validate:"required"`
	Amount float64 `json:"amount" validate:"required"`
}

type SendMoney struct {
	Receiver_id string  `json:"receiver_id"`
	Amount      float64 `json:"amount"`
}
