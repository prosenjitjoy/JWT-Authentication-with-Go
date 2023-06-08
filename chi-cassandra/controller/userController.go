package controller

import (
	"encoding/json"
	"main/helpers"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := helpers.CheckUserType(r.Context(), "ADMIN")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "unauthorized user"})
			return
		}

		usersMap, err := session.Query(`SELECT * from users`).Iter().SliceMap()

		if err != nil || len(usersMap) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to map slice"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(usersMap)
	}
}

func GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "user_id")

		err := helpers.MatchUserTypeToUID(r.Context(), userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": err.Error()})
			return
		}

		user, err := session.Query(`SELECT * FROM getUserByID WHERE user_id = ?`, userId).Iter().SliceMap()
		if err != nil || len(user) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to map slice"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
