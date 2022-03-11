package helpers

import (
	"github.com/Clementol/Wallet/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetWalletDetails(walletObj primitive.M, wallet models.Wallet) primitive.M {

	walletObj["user_id"] = wallet.User_id
	walletObj["wallet_id"] = wallet.Wallet_id
	walletObj["balance"] = wallet.Balance
	walletObj["created_at"] = wallet.Created_at
	walletObj["updated_at"] = wallet.Updated_at

	return walletObj

}
