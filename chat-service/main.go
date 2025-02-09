package main

import (
	"CNAD_Assignment_2/chat-service/routes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

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

	// 🚨 Debug: Print API Key Before Sending Request
	fmt.Println("🔑 Sending API Key:", openRouterAPIKey) // <-- ADD THIS

	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openRouterAPIKey). // ✅ Verify API key is correct
		SetHeader("Content-Type", "application/json").
		SetHeader("HTTP-Referer", "http://localhost:8081").
		SetHeader("X-Title", "Lion Befrienders Chatbot").
		SetBody(map[string]interface{}{
			"model": "google/gemini-2.0-flash-001",
			"messages": []map[string]string{
				{"role": "system", "content": "You are LionBee, a chatbot designed to help elderly users with fall-risk self-assessments."},
				{"role": "user", "content": req.Message},
			},
		}).
		Post("https://openrouter.ai/api/v1/chat/completions")

	if err != nil {
		fmt.Println("❌ Error sending request to OpenRouter:", err)
		http.Error(w, "Error communicating with OpenRouter API", http.StatusInternalServerError)
		return
	}

	fmt.Println("📩 Raw API Response:", string(resp.Body())) // ✅ Print full API response

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

func startChatbot() {
	router := mux.NewRouter()
	routes.SetupRoutes(router)

	// Apply CORS middleware
	handler := enableCORS(router)

	port := "8084"
	fmt.Println("🗨️ Chatbot microservice running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func startVoiceRecognition() {
	fmt.Println("🎤 Starting voice recognition...")

	cmd := exec.Command("python", "record_audio.py")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		fmt.Println("❌ Error running microphone recording:", err)
	} else {
		fmt.Println("✅ Microphone recording ran successfully")
	}
}

func main() {
	// ✅ Run Chatbot API & Voice Recognition in Parallel
	go startChatbot()       // 🔹 Runs chatbot in a separate goroutine
	startVoiceRecognition() // 🔹 Runs voice recording in the main goroutine
}
