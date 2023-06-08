package helpers

import (
	"main/model"
	"time"

	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v5"
)

var session *gocql.Session = model.DBSession()

var key string = "gojwt"

type SignedDetails struct {
	FirstName string
	LastName  string
	Email     string
	UserType  string
	UserID    string
	jwt.RegisteredClaims
}

func GenerateAllTokens(user model.User) (token string, refreshToken string, err error) {
	claims := &SignedDetails{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserID:    user.UserID,
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

func UpdateToken(new_token string, refresh_token string, user model.User) error {
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	err := session.Query(`UPDATE users SET new_token = ?, refresh_token = ?, updated_at = ? WHERE last_name = ? AND email = ? AND phone = ?`, new_token, refresh_token, updated_at, user.LastName, user.Email, user.Phone).Exec()
	return err
}

func ValidateToken(clientToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(clientToken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		msg = "failed to parse token"
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		msg = "the token is expired"
		return
	}

	return claims, ""
}
