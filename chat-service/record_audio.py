import sounddevice as sd
import vosk
import sys
import queue
import json

# ✅ Load Vosk Model
MODEL_PATH = "model"
model = vosk.Model(MODEL_PATH)

# ✅ Queue to store recorded audio
audio_queue = queue.Queue()

# ✅ Callback function to process audio chunks
def callback(indata, frames, time, status):
    if status:
        print(status, file=sys.stderr)
    audio_queue.put(bytes(indata))

# ✅ Configure recording settings
samplerate = 16000
blocksize = 8000

# ✅ Start recording
with sd.RawInputStream(samplerate=samplerate, blocksize=blocksize, dtype="int16",
                        channels=1, callback=callback):
    print("🎤 Recording... Speak into the microphone! (Press Ctrl+C to stop)")
    
    rec = vosk.KaldiRecognizer(model, samplerate)
    
    try:
        while True:
            data = audio_queue.get()
            if rec.AcceptWaveform(data):
                result = json.loads(rec.Result())
                text = result["text"]

                # ✅ Remove "the" at the beginning & end
                words = text.split()
                if words and words[0] == "the":
                    words.pop(0)  # Remove first word if it's "the"
                if words and words[-1] == "the":
                    words.pop()  # Remove last word if it's "the"

                cleaned_text = " ".join(words)
                print("You said:", cleaned_text)  # ✅ Prints cleaned text
    except KeyboardInterrupt:
        print("\n🛑 Stopping recording...")
