package main

import (
	"CNAD_Assignment_2/chat-service/routes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

// Request & Response Structures
type ChatRequest struct {
	Message string `json:"message"`
}

type ChatbotResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Hardcoded API Key (⚠️ Temporary solution)
var openRouterAPIKey = "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"

func chatHandler(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openRouterAPIKey). // ✅ Use API key
		SetHeader("Content-Type", "application/json").
		SetHeader("HTTP-Referer", "http://localhost:8081"). // ✅ Correct for local testing
		SetHeader("X-Title", "Lion Befrienders Chatbot").
		SetBody(map[string]interface{}{
			"model": "google/gemini-2.0-flash-001",
			"messages": []map[string]string{
				{"role": "system", "content": "You are LionBee, a friendly and patient chatbot designed to help elderly users perform a fall-risk self-assessment from home. Your goal is to provide clear, concise, and reassuring guidance while making the process as easy as possible. Keep instructions simple and step-by-step, avoiding medical jargon. If a user's responses indicate a high fall risk, gently suggest they seek medical advice or assistance. Prioritize empathy and clarity in all responses." +
					"When emphasizing important words or phrases, do NOT use Markdown formatting (such as # for headings or ** for bold). Instead, naturally highlight key terms by writing them in clear, readable sentences. For example, instead of **Important**, say: 'This is important: [text]'. Avoid using special characters that might confuse elderly users. Strictly no asterix involved."},
				{"role": "user", "content": req.Message},
			},
		}).
		Post("https://openrouter.ai/api/v1/chat/completions") // ✅ Correct API endpoint

	if err != nil {
		http.Error(w, "Error communicating with OpenRouter API", http.StatusInternalServerError)
		return
	}

	var chatbotResponse ChatbotResponse
	err = json.Unmarshal(resp.Body(), &chatbotResponse)
	if err != nil {
		http.Error(w, "Error parsing chatbot response", http.StatusInternalServerError)
		return
	}

	// Ensure chatbot response is not empty
	if len(chatbotResponse.Choices) == 0 || chatbotResponse.Choices[0].Message.Content == "" {
		http.Error(w, "Chatbot returned an empty response", http.StatusInternalServerError)
		return
	}

	// Send chatbot response back to frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatbotResponse)
}

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins for testing
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	router := mux.NewRouter()
	routes.SetupRoutes(router)

	// Apply CORS middleware
	handler := enableCORS(router)

	port := "8084"
	fmt.Println("Chatbot microservice running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
