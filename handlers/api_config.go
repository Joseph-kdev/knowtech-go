package handlers

import (
	"github.com/Joseph-kdev/knowtech-go/internal/db"
	"google.golang.org/genai"
)

type Apiconfig struct {
	DB     *db.Queries
	Gemini *genai.Client
}