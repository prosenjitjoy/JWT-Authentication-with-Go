package helpers

import (
	"fmt"
	"main/database"
	"main/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	UserType  string
	jwt.RegisteredClaims
}

var db *sqlx.DB = database.DBInstance()
var key string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(user models.UserInfo) (token string, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UID:       user.UserID,
		UserType:  user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(24))),
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(168))),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(key))
	if err != nil {
		panic(err)
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(key))
	if err != nil {
		panic(err)
	}

	return token, refreshToken, err
}

func UpdateAllTokens(token, refreshToken, userId string) {
	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	db.MustExec("UPDATE users SET token=$1, refresh_token=$2, updated_at=$3 WHERE user_id=$4", token, refreshToken, UpdatedAt, userId)
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		// msg = err.Error()
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		msg = "the token is expired"
		// msg = err.Error()
		return
	}
	fmt.Println("Claims", claims.UserType)
	return claims, ""
}
