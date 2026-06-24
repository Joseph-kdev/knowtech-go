package handlers

import (
	"net/http"

	"github.com/Joseph-kdev/knowtech-go/response"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	response.RespondWithJSON(w, 200, struct{}{})
}