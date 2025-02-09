package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

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
