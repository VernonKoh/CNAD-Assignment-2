import cv2
import numpy as np
import mediapipe as mp
from flask import Flask, request, send_file, jsonify, Response
import os
import time
import math

app = Flask(__name__)

mp_pose = mp.solutions.pose
pose = mp_pose.Pose(static_image_mode=True, min_detection_confidence=0.5, model_complexity=2)
pose2 = mp_pose.Pose(static_image_mode=False, min_detection_confidence=0.5, model_complexity=2)
mp_drawing = mp.solutions.drawing_utils

# Get the absolute path to the directory where your script is
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
UPLOAD_FOLDER = os.path.join(BASE_DIR, "uploads")  

if not os.path.exists(UPLOAD_FOLDER):
    os.makedirs(UPLOAD_FOLDER)

# Function to calculate angles between three points
def calculate_angle(a, b, c):
    ang = math.degrees(math.atan2(c[1] - b[1], c[0] - b[0]) -
                       math.atan2(a[1] - b[1], a[0] - b[0]))
    return abs(ang)

@app.route("/pose", methods=["POST"])
def process_image():
    file = request.files["image"]
    image = np.frombuffer(file.read(), np.uint8)
    image = cv2.imdecode(image, cv2.IMREAD_COLOR)

    image_rgb = cv2.cvtColor(image, cv2.COLOR_BGR2RGB)
    results = pose.process(image_rgb)

    if not results.pose_landmarks:
        return jsonify({"error": "No pose detected"})

    # ✅ Annotate image
    annotated_image = image.copy()
    mp_drawing.draw_landmarks(
        annotated_image, 
        results.pose_landmarks, 
        mp_pose.POSE_CONNECTIONS,
        landmark_drawing_spec=mp_drawing.DrawingSpec(color=(0, 255, 0), thickness=2, circle_radius=3),
        connection_drawing_spec=mp_drawing.DrawingSpec(color=(255, 0, 0), thickness=2)
    )

    # ✅ Extract keypoints for posture analysis
    landmarks = results.pose_landmarks.landmark
    def get_point(index):
        return (landmarks[index].x, landmarks[index].y)

    shoulder = get_point(11)  # Left Shoulder
    hip = get_point(23)       # Left Hip
    knee = get_point(25)      # Left Knee
    ankle = get_point(27)     # Left Ankle
    chin = get_point(1)       # Chin
    neck = get_point(0)       # Neck

    # ✅ Compute posture risk angles
    back_angle = calculate_angle(shoulder, hip, knee)
    head_angle = calculate_angle(neck, chin, shoulder)
    knee_angle = calculate_angle(hip, knee, ankle)

    risks = {}
    if back_angle < 160:
        risks["hunchback"] = "High Risk: Your back posture is bad."
    if head_angle < 50:
        risks["forward_head"] = "Warning: Forward Head Posture detected."
    if knee_angle < 160:
        risks["bent_knees"] = "Potential fall risk: Knees are bent."

    # ✅ Save the processed image
    output_path = os.path.join(UPLOAD_FOLDER, "output.jpg")
    cv2.imwrite(output_path, annotated_image)

    return jsonify({
        "image_url": "/processed-image",
        "analysis": risks
    })

@app.route("/processed-image", methods=["GET"])
def get_processed_image():
    output_path = os.path.join(UPLOAD_FOLDER, "output.jpg")
    return send_file(output_path, mimetype="image/jpeg")

# ✅ Real-Time Webcam Stream
def generate_frames():
    cap = cv2.VideoCapture(0)  # Use webcam (change to 1 if using external camera)
    
    while cap.isOpened():
        success, frame = cap.read()
        if not success:
            break

        # Convert frame to RGB
        image_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
        results = pose2.process(image_rgb)

        if results.pose_landmarks:
            # Draw pose landmarks
            mp_drawing.draw_landmarks(
                frame, results.pose_landmarks, mp_pose.POSE_CONNECTIONS,
                landmark_drawing_spec=mp_drawing.DrawingSpec(color=(0, 255, 0), thickness=2, circle_radius=3),
                connection_drawing_spec=mp_drawing.DrawingSpec(color=(255, 0, 0), thickness=2)
            )

            # Extract keypoints
            landmarks = results.pose_landmarks.landmark
            def get_point(index):
                return (landmarks[index].x, landmarks[index].y)

            shoulder = get_point(11)  # Left Shoulder
            hip = get_point(23)       # Left Hip
            knee = get_point(25)      # Left Knee
            ankle = get_point(27)     # Left Ankle
            chin = get_point(1)       # Chin
            neck = get_point(0)       # Neck

            # ✅ Compute posture angles
            back_angle = calculate_angle(shoulder, hip, knee)
            head_angle = calculate_angle(neck, chin, shoulder)
            knee_angle = calculate_angle(hip, knee, ankle)

            # ✅ Detect risky postures
            risk_text = ""
            if back_angle < 160:
                risk_text = "Hunchback detected!"
            elif head_angle < 50:
                risk_text = "Forward Head Posture!"
            elif knee_angle < 160:
                risk_text = "Bent Knees: Fall Risk!"

            # Draw warning text on frame
            if risk_text:
                cv2.putText(frame, risk_text, (50, 50), cv2.FONT_HERSHEY_SIMPLEX, 
                            1, (0, 0, 255), 2, cv2.LINE_AA)

        # Encode frame for streaming
        _, buffer = cv2.imencode('.jpg', frame)
        frame_bytes = buffer.tobytes()
        
        yield (b'--frame\r\n'
               b'Content-Type: image/jpeg\r\n\r\n' + frame_bytes + b'\r\n')

    cap.release()

@app.route("/webcam")
def webcam_feed():
    return Response(generate_frames(), mimetype='multipart/x-mixed-replace; boundary=frame')

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
