package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Joseph-kdev/knowtech-go/response"
	"google.golang.org/genai"
)

type chatMessage struct {
	Role string `json:"role"`
	Parts string `json:"parts"`
}

type chatRequest struct {
	History []chatMessage `json:"history"`
	Message string `json:"message"`
}

type chatResponse struct {
	Reply string `json:"reply"`
	History []chatMessage `json:"history"`
} 

func (apiCfg *Apiconfig) HandleChats(w http.ResponseWriter, r *http.Request) {
	type parameters chatRequest

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		response.RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	history := mapHistoryToContent(params.History)

	systemText := os.Getenv("SYSTEM_PROMPT")
	if systemText == "" {
		response.RespondWithError(w, 500, "error getting system prompt")
		return
	}

	systemPrompt := &genai.Content{
		Parts: []*genai.Part{
			genai.NewPartFromText(systemText),
		},
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: systemPrompt,
	}

	chat, err := apiCfg.Gemini.Chats.Create(r.Context(), "gemini-2.5-flash", config, history)
	if err != nil {
		response.RespondWithError(w, 500, fmt.Sprintf("Error starting chat: %v", err))
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		response.RespondWithError(w, 500, "Error: streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	var fullReply strings.Builder

	for chunk, err := range chat.SendMessageStream(r.Context(), genai.Part{Text: params.Message}) {
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}

		text := chunk.Text()
		fullReply.WriteString(text)

		payload, _ := json.Marshal(map[string]string{"chunk": text})
		fmt.Fprintf(w, "data: %s\n\n", payload)
		flusher.Flush()
	}

	updatedHistory := append(params.History,
		chatMessage{Role: "user", Parts: params.Message},
		chatMessage{Role: "model", Parts: fullReply.String()},
	)

	finalPayload, _ := json.Marshal(map[string]any{
		"done": true,
		"history": updatedHistory,
	})

	fmt.Fprintf(w, "event: done\ndata: %s\n\n", finalPayload)
	flusher.Flush()
}

func mapHistoryToContent(history []chatMessage) []*genai.Content {
	contents := make([]*genai.Content, 0, len(history))

	for _, m := range history {
		contents = append(contents, &genai.Content{
			Role: m.Role,
			Parts: []*genai.Part{
				genai.NewPartFromText(m.Parts),
			},
		})
	}

	return contents
}