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
		fmt.Println("❌ Failed to parse request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("📥 Received message from user:", req.Message)
	fmt.Println("🔑 Using API Key:", openRouterAPIKey) // Debug API Key

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openRouterAPIKey).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"model": "deepseek/deepseek-chat",
			"messages": []map[string]string{
				{"role": "user", "content": req.Message},
			},
		}).
		Post("https://openrouter.ai/api/v1/chat/completions")

	if err != nil {
		fmt.Println("❌ Error sending request to OpenRouter:", err)
		http.Error(w, "Error communicating with OpenRouter API", http.StatusInternalServerError)
		return
	}

	// 🚨 Print full API response
	fmt.Println("📩 Raw API Response:", string(resp.Body()))
	fmt.Println("API Status Code:", resp.StatusCode())

	if resp.StatusCode() != 200 {
		http.Error(w, fmt.Sprintf("API error: %d", resp.StatusCode()), http.StatusInternalServerError)
		return
	}

	var chatbotResponse ChatbotResponse
	err = json.Unmarshal(resp.Body(), &chatbotResponse)
	if err != nil {
		fmt.Println("❌ Failed to parse API response:", err)
		http.Error(w, "Error parsing chatbot response", http.StatusInternalServerError)
		return
	}

	if len(chatbotResponse.Choices) == 0 || chatbotResponse.Choices[0].Message.Content == "" {
		fmt.Println("⚠️ API returned empty response")
		http.Error(w, "Chatbot returned an empty response", http.StatusInternalServerError)
		return
	}

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

	// Register the handler for the /chat route
	router.HandleFunc("/chat", chatHandler).Methods("POST")

	routes.SetupRoutes(router)

	// Apply CORS middleware
	handler := enableCORS(router)

	port := "8084"
	fmt.Println("Chatbot microservice running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
