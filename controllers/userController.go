package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/Clementol/Wallet/database"
	helper "github.com/Clementol/Wallet/helpers"
	"github.com/Clementol/Wallet/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection = db.OpenCollection(db.Client, "user")
var walletCollection = db.OpenCollection(db.Client, "wallet")
var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var wallet models.Wallet

		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate the data based on user struct

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			msg := "error occurred while checking for email"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		password := helper.HashPassword(user.Password)
		user.Password = password

		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			msg := "error occured while checking for phone number"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		if emailCount > 0 || phoneCount > 0 {
			msg := "this email or phone number already exist"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		user.Active = true
		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.First_Name, user.Last_Name, user.User_id, user.Active)
		user.Token = token
		user.Refresh_Token = refreshToken

		result, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "user was not registered"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Create Wallet for user
		wallet.ID = primitive.NewObjectID()
		wallet.Wallet_id = wallet.ID.Hex()
		wallet.User_id = user.User_id
		wallet.Balance = 0.00
		wallet.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		wallet.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		_, walletErr := walletCollection.InsertOne(ctx, wallet)

		if walletErr != nil {
			msg := "Unable to create Wallet"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		foundUserObj := bson.M{}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": &user.Email}).Decode(&foundUser)
		if err != nil {
			msg := "invalid user credentials"
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		defer cancel()
		passwordInvalid, msg := helper.VerifyPassword(user.Password, foundUser.Password)
		if !passwordInvalid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}
		defer cancel()

		token, refreshToken, _ := helper.GenerateAllTokens(foundUser.Email, foundUser.First_Name, foundUser.Last_Name, foundUser.User_id, foundUser.Active)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		userDetails := helper.GetUserDetails(foundUserObj, foundUser)

		c.JSON(http.StatusOK, userDetails)
	}
}
