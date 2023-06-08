package controller

import (
	"encoding/json"
	"main/helpers"
	"main/model"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

type status map[string]interface{}

var session *gocql.Session = model.DBSession()
var validate = validator.New()

func Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		err = validate.Struct(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "failed to validate json:"})
			return
		}

		var emailCount int
		err = session.Query(`SELECT COUNT(*) FROM chijwt.users WHERE last_name = ? AND email = ?`, user.LastName, user.Email).Scan(&emailCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to count email:"})
			return
		}

		var phoneCount int
		err = session.Query(`SELECT COUNT(*) FROM users WHERE last_name = ? AND email = ? AND phone = ?`, user.LastName, user.Email, user.Phone).Scan(&phoneCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to count phone:"})
			return
		}

		if emailCount > 0 || phoneCount > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "duplicate email or phone:"})
			return
		}

		user.Password = helpers.HashPassword(user.Password)
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UserID = uuid.NewString()
		// fmt.Println(user)

		new_token, refresh_token, err := helpers.GenerateAllTokens(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to generate token:"})
			return
		}

		user.NewToken = new_token
		user.RefreshToken = refresh_token

		err = session.Query(`INSERT INTO users (first_name, last_name, email, password, phone, new_token, user_type, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.NewToken, user.UserType, user.UserID, user.CreatedAt, user.UpdatedAt).Exec()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to insert data:"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status{"InsertedID": user.UserID})
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		if user.LastName == "" || user.Email == "" || user.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "enter last_name, email and phone"})
			return
		}

		usersMap, err := session.Query(`SELECT * FROM users WHERE last_name = ? AND email = ?`, user.LastName, user.Email).Iter().SliceMap()
		if err != nil || len(usersMap) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to map slice"})
			return
		}

		jsonStr, err := json.Marshal(usersMap[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "json marshling failed"})
			return
		}

		var foundUser model.User
		if err := json.Unmarshal(jsonStr, &foundUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "json unmarshling failed"})
			return
		}

		passwordIsValid, msg := helpers.VerifyPassword(user.Password, foundUser.Password)

		if passwordIsValid != true {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": msg})
			return
		}

		newToken, refToken, err := helpers.GenerateAllTokens(foundUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to generate token:"})
			return
		}

		err = helpers.UpdateToken(newToken, refToken, foundUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to update token:" + err.Error()})
			return
		}

		usersMap, err = session.Query(`SELECT * FROM getUserByID WHERE user_id = ?`, foundUser.UserID).Iter().SliceMap()
		if err != nil || len(usersMap) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to map slice"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(usersMap[0])
	}
}
