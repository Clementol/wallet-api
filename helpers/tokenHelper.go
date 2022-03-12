package helpers

import (
	"context"
	"log"
	"os"
	"time"

	db "github.com/Clementol/Wallet/database"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	Active     bool
	jwt.StandardClaims
}

var userCollection *mongo.Collection = db.OpenCollection(db.Client, "user")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, uid string, active bool) (signedToken, signedRefreshToken string, err error) {

	claims := &SignedDetails{
		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		Uid:        uid,
		Active:     active,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	var tokens, err1 = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err1 != nil {
		log.Fatal(err)
		return
	}

	return tokens, refreshToken, err
}

func UpdateAllTokens(signedToken, signedRefreshToken, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	// upsert := true
	filter := bson.M{"user_id": userId}
	// opt := options.UpdateOptions{
	// 	Upsert: &upsert,
	// }
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{primitive.E{Key: "$set", Value: updateObj}},
		// &opt,
	)

	defer cancel()

	if err != nil {
		log.Fatal(err)
		return
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = "Error while Validating " + err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Authentication Error" + err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Please login in Again" + err.Error()
		return
	}

	return claims, msg
}
