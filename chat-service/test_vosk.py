import vosk
import sys
import wave
import json

# Load Vosk Model
model = vosk.Model("model")  # Ensure the "model" folder exists

# Open the audio file
audio_path = sys.argv[1]
wf = wave.open(audio_path, "rb")
rec = vosk.KaldiRecognizer(model, wf.getframerate())

# Process audio and print transcript
while True:
    data = wf.readframes(4000)
    if len(data) == 0:
        break
    if rec.AcceptWaveform(data):
        result = json.loads(rec.Result())
        print(result["text"])  # ✅ Prints the converted text

# Final output
final_result = json.loads(rec.FinalResult())
print(final_result["text"])
