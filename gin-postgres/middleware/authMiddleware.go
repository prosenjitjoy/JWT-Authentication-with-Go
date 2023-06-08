package middleware

import (
	"fmt"
	"main/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "no authorization header provided!"})
			ctx.Abort()
			return
		}

		fmt.Println(clientToken, "$$$")
		claims, msg := helpers.ValidateToken(clientToken)
		if msg != "" {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": msg})
			ctx.Abort()
			return
		}

		fmt.Println(claims.FirstName, claims.ExpiresAt, claims.UserType, "$$$$")

		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.FirstName)
		ctx.Set("last_name", claims.LastName)
		ctx.Set("uid", claims.UID)
		ctx.Set("user_type", claims.UserType)
		fmt.Println(claims.UserType, "*")
	}
}
