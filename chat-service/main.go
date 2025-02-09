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
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+openRouterAPIKey).
		SetHeader("Content-Type", "application/json").
		SetHeader("HTTP-Referer", "http://localhost:8081"). // ✅ Corrected for local testing
		SetHeader("X-Title", "Lion Befrienders Chatbot").
		SetBody(map[string]interface{}{
			"model": "deepseek/deepseek-chat",
			"messages": []map[string]string{
				{"role": "system", "content": "You are a fall-risk self-assessment assistant for elderly users. Provide step-by-step guidance in a simple way."},
				{"role": "user", "content": req.Message},
			},
		}).
		Post("https://openrouter.ai/api/v1/chat/completions")

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
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081") // Allow frontend
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS request)
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

	// ✅ Use Python to run microphone recording
	cmd := exec.Command("python", "record_audio.py")
	cmd.Stderr = os.Stderr // Print errors to terminal
	cmd.Stdout = os.Stdout // Print output to terminal
	err := cmd.Run()

	if err != nil {
		fmt.Println("Error running microphone recording:", err)
	} else {
		fmt.Println("Microphone recording ran successfully")
	}

	// // ✅ Test with a sample audio file (Replace with your own)
	// audioFile := "knees.wav"

	// // Run the Python script to process the audio
	// cmd := exec.Command("python", "test_vosk.py", audioFile)
	// output, err := cmd.CombinedOutput()

	// if err != nil {
	// 	fmt.Println("Error running Vosk:", err)
	// } else {
	// 	fmt.Println("Transcribed Text:\n", string(output))
	// }
}
