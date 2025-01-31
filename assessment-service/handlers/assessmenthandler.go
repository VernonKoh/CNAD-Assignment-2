package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	mediapipeAPIURL   = "http://127.0.0.1:5000/pose"
	mediapipeImageURL = "http://127.0.0.1:5000/processed-image"
)

// UploadHandler handles file uploads and forwards them to the MediaPipe API
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save temporarily
	filePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	defer os.Remove(filePath)

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Forward the file to MediaPipe API
	response, err := sendToMediaPipeAPI(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to call MediaPipe API: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// ✅ Read JSON response from Python
	var analysisData map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&analysisData)
	if err != nil {
		http.Error(w, "Failed to read posture analysis data", http.StatusInternalServerError)
		return
	}

	// ✅ Fetch processed image from Python API
	imageResp, err := http.Get(mediapipeImageURL)
	if err != nil {
		http.Error(w, "Failed to retrieve processed image", http.StatusInternalServerError)
		return
	}
	defer imageResp.Body.Close()

	// ✅ Set response headers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"analysis":  analysisData["analysis"],
		"image_url": mediapipeImageURL,
	})
}

// sendToMediaPipeAPI forwards the file to the Python MediaPipe API
func sendToMediaPipeAPI(filePath string) (*http.Response, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("error creating form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("error copying file content to form: %v", err)
	}
	writer.Close()

	// Send the request to the MediaPipe API
	req, err := http.NewRequest("POST", mediapipeAPIURL, body)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	return client.Do(req)
}
