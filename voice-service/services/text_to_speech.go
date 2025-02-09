package services

import (
	"fmt"
	"os"
	"os/exec"
)

// TextToSpeech - Converts text to speech using Coqui TTS
func TextToSpeech(text string) []byte {
	tmpFile, err := os.CreateTemp("", "speech-*.wav")
	if err != nil {
		fmt.Println("Error creating temp file:", err)
		return nil
	}
	defer os.Remove(tmpFile.Name())

	// Run Coqui TTS CLI
	cmd := exec.Command("tts", "--text", text, "--out_path", tmpFile.Name())
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running TTS:", err)
		return nil
	}

	audioData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		fmt.Println("Error reading TTS output:", err)
		return nil
	}

	return audioData
}
