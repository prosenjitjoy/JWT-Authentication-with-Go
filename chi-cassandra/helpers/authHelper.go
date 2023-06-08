package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func VerifyPassword(providedPassword, userPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}

	return check, msg
}

func CheckUserType(c context.Context, role string) error {
	userType := c.Value("user_type")
	if userType != role {
		return fmt.Errorf("unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypeToUID(c context.Context, user_id string) error {
	userType := c.Value("user_type")
	userId := c.Value("user_id")

	if userType == "USER" && userId != user_id {
		return fmt.Errorf("unauthorized to access this resource")
	}
	return nil
}
