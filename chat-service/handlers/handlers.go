package handlers

import (
	"CNAD_Assignment_2/chat-service/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Your OpenRouter API Key
var openRouterAPIKey = "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"

// ChatHandler processes chatbot requests
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("❌ Error decoding request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Println("📥 Received message from user:", req.Message)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openRouterAPIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "deepseek/deepseek-chat",
			"messages": []map[string]string{
				{"role": "system", "content": "You are a helpful assistant for elderly users."},
				{"role": "user", "content": req.Message},
			},
		}).
		Post("https://openrouter.ai/api/v1/chat/completions")

	if err != nil {
		log.Println("❌ Error calling OpenRouter API:", err)
		http.Error(w, "Error communicating with OpenRouter API", http.StatusInternalServerError)
		return
	}

	log.Println("✅ API Response:", resp.String())

	var chatbotResponse models.ChatbotResponse
	err = json.Unmarshal(resp.Body(), &chatbotResponse)
	if err != nil {
		log.Println("❌ Error parsing API response:", err)
		http.Error(w, "Error parsing chatbot response", http.StatusInternalServerError)
		return
	}

	// Check if DeepSeek returned an actual message
	if len(chatbotResponse.Choices) == 0 || chatbotResponse.Choices[0].Message.Content == "" {
		log.Println("⚠️ API returned empty response")
		http.Error(w, "Chatbot returned an empty response", http.StatusInternalServerError)
		return
	}

	// Send chatbot response to frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatbotResponse)
	log.Println("📤 Sent response to frontend:", chatbotResponse.Choices[0].Message.Content)
}
