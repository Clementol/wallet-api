package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/Clementol/Wallet/database"
	helper "github.com/Clementol/Wallet/helpers"
	"github.com/Clementol/Wallet/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var historyCollection = db.OpenCollection(db.Client, "history")

func FundWallet() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var foundUser models.User
		var userWallet models.Wallet
		var foundWallet models.Wallet
		var fundWallet models.FundWallet
		// var history models.History
		updateObj := bson.M{}

		userId := c.MustGet("user_id").(string)

		if err := c.Bind(&fundWallet); err != nil {
			log.Fatal(err.Error())
		}

		validationErr := validate.Struct(&fundWallet)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			msg := "user not found"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if !foundUser.Active {
			msg := "account currently deactivated"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		err = walletCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundWallet)
		if err != nil {
			msg := "wallet not found"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		updatedBalance := foundWallet.Balance + fundWallet.Amount
		updateObj["balance"] = updatedBalance
		updateObj["updated_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.M{"user_id": userId}
		updateWallet := bson.M{"$set": updateObj}
		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err = walletCollection.FindOneAndUpdate(ctx,
			filter, updateWallet, opt).Decode(&userWallet)

		if err != nil {
			msg := "Can't fund wallet" + err.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		updatedWalletObj := bson.M{}

		walletDetails := helper.GetWalletDetails(updatedWalletObj, userWallet)

		c.JSON(http.StatusAccepted, walletDetails)

	}
}

func SendMoney() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var sendMoney models.SendMoney
		var foundWallet models.Wallet
		var userWallet models.Wallet
		var foundUser models.User
		// var history models.History

		var foundReceiverWallet models.Wallet
		var receiverWallet models.Wallet
		var foundReceiver models.User

		if err := c.Bind(&sendMoney); err != nil {
			log.Fatal(err.Error())
		}

		userId := c.MustGet("user_id").(string)
		// senderName := c.MustGet("first_name").(string)

		if sendMoney.Amount <= 0 {
			msg := "amount must be greater than zero"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			msg := "user not found"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		if !foundUser.Active {
			msg := "account currently deactivated"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		}

		err = walletCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&foundWallet)
		if err != nil {
			msg := "wallet not found"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if foundWallet.Balance == 0 {
			msg := "Please fund your wallet"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if foundWallet.Balance < sendMoney.Amount {
			msg := fmt.Sprintf("insufficient balance of %v", foundWallet.Balance)
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		// find receiver
		err = userCollection.FindOne(ctx, bson.M{"user_id": sendMoney.Receiver_id}).Decode(&foundReceiver)
		defer cancel()
		if err != nil {
			msg := "receiver not found"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		// debit sender wallet
		remainingBalance := foundWallet.Balance - sendMoney.Amount
		userWaletObj := bson.M{}

		userWaletObj["balance"] = remainingBalance
		userWaletObj["updated_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filter := bson.M{"user_id": userId}
		updateWallet := bson.M{"$set": userWaletObj}
		opt := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err = walletCollection.FindOneAndUpdate(ctx,
			filter, updateWallet, opt).Decode(&userWallet)

		if err != nil {
			msg := "Transaction not completed" + err.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		// credit receiver wallet
		err = walletCollection.FindOne(ctx, bson.M{"user_id": sendMoney.Receiver_id}).Decode(&foundReceiverWallet)
		defer cancel()
		if err != nil {
			msg := "wallet not found"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		receiverWaletObj := bson.M{}
		updatedReceiverBalance := foundReceiverWallet.Balance + sendMoney.Amount

		receiverWaletObj["balance"] = updatedReceiverBalance
		receiverWaletObj["updated_at"], _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		filterReciver := bson.M{"user_id": sendMoney.Receiver_id}
		updateReceiverWallet := bson.M{"$set": receiverWaletObj}
		optRec := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err = walletCollection.FindOneAndUpdate(ctx,
			filterReciver, updateReceiverWallet,
			optRec).Decode(&receiverWallet)

		if err != nil {
			msg := "Transaction not completed" + err.Error()
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusAccepted, receiverWallet)

	}
}
