package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	fmt.Println("State 3", userType)
	if userType != role {
		return fmt.Errorf("unauthorized to access this resource")
	}
	return nil
}

func MatchUserTypeToUid(c *gin.Context, userId string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")

	if userType == "USER" && uid != userId {
		return fmt.Errorf("unauthorized to access this resource")
	}

	return CheckUserType(c, userType)
}
