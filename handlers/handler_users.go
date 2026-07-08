package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"github.com/Joseph-kdev/knowtech-go/response"
)

func (apiCfg *Apiconfig) AddUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserId string `json:"user_id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}
	
	user, err := apiCfg.DB.AddUserToDatabase(r.Context(), db.AddUserToDatabaseParams{
		ID: params.UserId,
		Email: params.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error adding user: %v", err))
		return
	}

	response.RespondWithJSON(w, 200, user)
}