package handler

import (
	"chirpy/internal/database"
	"chirpy/reponse"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	user := params{}
	err := decoder.Decode(&user)
	if err != nil {
		reponse.RespondWithError(w, 422, "Invalid email format")
		// respondWithError

		return
	}
	userDb, err := cfg.Db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Email:     user.Email,
	})
	if err != nil {
		reponse.RespondWithError(w, 500, "user not created")
		fmt.Println(err)
		return
	}
	userResult := User{
		ID:        userDb.ID,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
		Email:     userDb.Email,
	}
	// data, err := json.Marshal(userDb)
	// if err != nil {
	// 	respondWithError(w, 500, "user not created")
	// 	return
	// }
	reponse.RespondWithJSON(w, 201, userResult)
}
