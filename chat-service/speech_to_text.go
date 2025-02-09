package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// SpeechToText - Uses Python Vosk for Speech Recognition
func SpeechToText(audioData []byte) string {
	tmpFile, err := os.CreateTemp("", "audio-*.wav")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return ""
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(audioData)
	if err != nil {
		fmt.Println("Error writing to temp file:", err)
		return ""
	}
	tmpFile.Close()

	// ✅ Run Python script using full path
	cmd := exec.Command("python", "chat-service/test_vosk.py", tmpFile.Name())
	var output bytes.Buffer
	cmd.Stdout = &output
	err = cmd.Run()

	if err != nil {
		fmt.Println("Error running Vosk:", err)
		return ""
	}

	return output.String()
}
