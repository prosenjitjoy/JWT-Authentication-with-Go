package controller

import (
	"main/database"
	"main/helpers"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB = database.DBInstance()
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func VerifyPassword(providedPassword, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	return err != nil
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.UserInfo

		if err := ctx.BindJSON(&user); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(user)
		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var emailCount int
		err = db.Get(&emailCount, "SELECT COUNT(*) FROM users WHERE email=$1", user.Email)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Password = HashPassword(user.Password)

		var phoneCount int
		err = db.Get(&phoneCount, "SELECT COUNT(*) FROM users WHERE phone=$1", user.Phone)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if emailCount > 0 || phoneCount > 0 {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UserID = uuid.NewString()

		token, refreshToken, err := helpers.GenerateAllTokens(user)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		user.Token = token
		user.RefreshToken = refreshToken

		_, err = db.NamedExec("INSERT INTO users (first_name, last_name, password, email, phone, token, user_type, refresh_token, created_at, updated_at, user_id) VALUES (:first_name, :last_name, :password, :email, :phone, :token, :user_type, :refresh_token, :created_at, :updated_at, :user_id)", user)

		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(http.StatusOK, gin.H{"InsertedID": user.UserID})
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.UserInfo

		if err := ctx.BindJSON(&user); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var foundUser models.User
		err := db.Get(&foundUser, "SELECT * FROM users WHERE email=$1", user.Email)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		foundUserInfo := foundUser.UserInfo
		passwordIsValid := VerifyPassword(user.Password, foundUserInfo.Password)

		if !passwordIsValid {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, refreshToken, err := helpers.GenerateAllTokens(foundUserInfo)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		helpers.UpdateAllTokens(token, refreshToken, foundUserInfo.UserID)

		err = db.Get(&foundUser, "SELECT * FROM users WHERE user_id=$1", foundUserInfo.UserID)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		users := []models.User{}
		err := db.Select(&users, "SELECT * FROM users")
		if err != nil {
			panic(err)
		}

		ctx.IndentedJSON(http.StatusOK, users)
	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		err := helpers.MatchUserTypeToUid(ctx, userId)

		if err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err = db.Get(&user, "SELECT * FROM users WHERE user_id=$1", userId)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.IndentedJSON(http.StatusOK, user)
	}
}
