package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/tonge3199/go-RSS-project/internal/database"
	"net/http"
	"time"
)

func (config *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	/* user input came from r.body
	   w : write HTTP Responses.
	*/
	type parameters struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params) // read in JSON data from r.Body.-> &params,   隐式赋值
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := config.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (config *apiConfig) handlerUserGet(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
