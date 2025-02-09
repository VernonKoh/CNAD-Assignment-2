package models

// ChatRequest represents the request body for chatbot requests
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatbotResponse represents the response from the chatbot API
type ChatbotResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
