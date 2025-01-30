import cv2
import numpy as np
import mediapipe as mp
from flask import Flask, request, send_file, jsonify
import os

app = Flask(__name__)

mp_pose = mp.solutions.pose
pose = mp_pose.Pose(static_image_mode=True, min_detection_confidence=0.5, model_complexity=2)
mp_drawing = mp.solutions.drawing_utils

# Get the absolute path to the directory where your script is
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
UPLOAD_FOLDER = os.path.join(BASE_DIR, "uploads")  # construct absolute path

if not os.path.exists(UPLOAD_FOLDER):
    os.makedirs(UPLOAD_FOLDER)

@app.route("/pose", methods=["POST"])
def process_image():
    file = request.files["image"]
    image = np.frombuffer(file.read(), np.uint8)
    image = cv2.imdecode(image, cv2.IMREAD_COLOR)

    image_rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    results = pose.process(image_rgb)

    if not results.pose_landmarks:
        return jsonify({"error": "No pose detected"})

    # ✅ Force drawing of landmarks
    annotated_image = image.copy()
    mp_drawing.draw_landmarks(
        annotated_image, 
        results.pose_landmarks, 
        mp_pose.POSE_CONNECTIONS,
        landmark_drawing_spec=mp_drawing.DrawingSpec(color=(0, 255, 0), thickness=2, circle_radius=3),  # Green landmarks
        connection_drawing_spec=mp_drawing.DrawingSpec(color=(255, 0, 0), thickness=2)  # Red connections
    )

    print(f"Shape of annotated image: {annotated_image.shape}")  # debugging

    # ✅ Save and return the processed image
    output_path = os.path.join(UPLOAD_FOLDER, "output.jpg")
    print(f"Saving to: {output_path}") # debugging
    cv2.imwrite(output_path, annotated_image)

    return send_file(output_path, mimetype="image/jpeg")

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)