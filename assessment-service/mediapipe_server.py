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
pose2 = mp_pose.Pose(min_detection_confidence=0.5, min_tracking_confidence=0.5)
mp_drawing = mp.solutions.drawing_utils

# Get the absolute path to the directory where your script is
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
UPLOAD_FOLDER = os.path.join(BASE_DIR, "uploads")  
PROCESSED_FOLDER = os.path.join(BASE_DIR, "processed") 
VIDEO_UPLOAD_FOLDER = os.path.join(BASE_DIR, "uploads")  
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER
app.config['PROCESSED_FOLDER'] = PROCESSED_FOLDER
app.config['VIDEO_UPLOAD_FOLDER'] = VIDEO_UPLOAD_FOLDER

if not os.path.exists(UPLOAD_FOLDER):
    os.makedirs(UPLOAD_FOLDER)

if not os.path.exists(PROCESSED_FOLDER):
    os.makedirs(PROCESSED_FOLDER)

if not os.path.exists(VIDEO_UPLOAD_FOLDER):
    os.makedirs(VIDEO_UPLOAD_FOLDER)


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


# Function to calculate angle between three points
def calculate_gait_angle(a, b, c):
    """
    Calculate the angle between three points (a, b, c) in degrees.
    Each point is a tuple (x, y).
    """
    import math
    # Convert points to vectors
    vector1 = (a[0] - b[0], a[1] - b[1])
    vector2 = (c[0] - b[0], c[1] - b[1])
    # Compute dot product
    dot_product = vector1[0] * vector2[0] + vector1[1] * vector2[1]
    # Compute magnitudes
    magnitude1 = math.sqrt(vector1[0]**2 + vector1[1]**2)
    magnitude2 = math.sqrt(vector2[0]**2 + vector2[1]**2)
    # Compute cosine of angle
    cosine_angle = dot_product / (magnitude1 * magnitude2)
    # Clamp cosine_angle to avoid numerical errors
    cosine_angle = max(-1.0, min(1.0, cosine_angle))
    # Compute angle in degrees
    angle = math.degrees(math.acos(cosine_angle))
    return angle

# Process a single frame for gait analysis
risk_text_buffer = []  # Store the risk text for multiple frames

def process_frame(frame, pose):
    global risk_text_buffer

    # Convert frame to RGB
    image_rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)
    results = pose.process(image_rgb)
    
    risk_text = ""  # Text to draw on the frame
    
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
            """Helper function to extract x, y coordinates of a landmark."""
            return (landmarks[index].x, landmarks[index].y)
        
        # Define key points for gait analysis
        left_hip = get_point(23)       # Left Hip
        left_knee = get_point(25)     # Left Knee
        left_ankle = get_point(27)    # Left Ankle
        right_hip = get_point(24)     # Right Hip
        right_knee = get_point(26)    # Right Knee
        right_ankle = get_point(28)   # Right Ankle
        
        # Compute gait-specific angles
        left_knee_angle = calculate_gait_angle(left_hip, left_knee, left_ankle)
        right_knee_angle = calculate_gait_angle(right_hip, right_knee, right_ankle)
        
        # Detect gait abnormalities
        if left_knee_angle < 160 or right_knee_angle < 160:
            risk_text = "Abnormal Knee Angle: Risk of Unstable Gait!"
        elif abs(left_knee_angle - right_knee_angle) > 20:
            risk_text = "Asymmetric Knee Angles: Risk of Uneven Gait!"
        
        # Update the risk text buffer
        if risk_text:
            risk_text_buffer.append(risk_text)  # Add the new risk text
            if len(risk_text_buffer) > 60:  # Limit the buffer size to 60 frames (adjust as needed)
                risk_text_buffer.pop(0)  # Remove the oldest entry
    
    # Draw risk text on the frame (if there's any in the buffer)
    if risk_text_buffer:
        text_to_display = risk_text_buffer[0]  # Display the most recent risk text
        cv2.putText(frame, text_to_display, (50, 50), cv2.FONT_HERSHEY_SIMPLEX, 
                    1.5, (0, 0, 255), 3, cv2.LINE_AA)  # Increased size and thickness

    return frame

      
def process_video(input_path, output_path):
    cap = cv2.VideoCapture(input_path)
    if not cap.isOpened():
        print(f"Error: Could not open video file at {input_path}")
        return  # Exit if video cannot be opened
    
    fps = int(cap.get(cv2.CAP_PROP_FPS))
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
    fourcc = cv2.VideoWriter_fourcc(*'H264')
    out = cv2.VideoWriter(output_path, fourcc, fps, (width, height))
    
    if not out.isOpened():
        print(f"Error: Could not open output video for writing at {output_path}")
        cap.release()
        return  # Exit if output video cannot be opened
    
    while cap.isOpened():
        success, frame = cap.read()
        if not success:
            break
        # Process each frame
        processed_frame = process_frame(frame, pose2)
        out.write(processed_frame)
    
    cap.release()
    out.release()
    cv2.destroyAllWindows()
    print("Video processing complete: Released resources.")

@app.route("/upload_video", methods=["POST"])
def upload_video():
    # Check if a file was uploaded
    if "file" not in request.files:
        return jsonify({"error": "No file part"}), 400
    file = request.files["file"]
    if file.filename == "":
        return jsonify({"error": "No selected file"}), 400
    
    # Save the uploaded file with its original name
    filename = file.filename  # Use the original filename
    filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
    try:
        file.save(filepath)
    except Exception as e:
        print(f"Error saving file: {e}")
        return jsonify({"error": "Failed to save file"}), 500
    
    # Process the video and save it as "processed.mp4"
    output_filename = "processed.mp4"  # Fixed name for the processed video
    output_path = os.path.join(app.config['PROCESSED_FOLDER'], output_filename)
    process_video(filepath, output_path)
    
    try:
        os.remove(filepath)  # Remove the original video after processing
    except Exception as e:
        print(f"Error deleting file: {e}")
        return jsonify({"error": "Video Processing Complete, but original video could not be deleted."}), 200
    
    # Return the path to the processed video
    return jsonify({"message": "Video Processing Complete", "download_url": f"/download/{output_filename}"}), 200

    

    
# Route to download processed video
@app.route("/download/<filename>")
def download_video(filename):
    return send_file(os.path.join(app.config['PROCESSED_FOLDER'], filename), as_attachment=True)


#webcam
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
                """Helper function to extract x, y coordinates of a landmark."""
                return (landmarks[index].x, landmarks[index].y)
            
            # Define key points for gait analysis
            left_hip = get_point(23)       # Left Hip
            left_knee = get_point(25)     # Left Knee
            left_ankle = get_point(27)    # Left Ankle
            right_hip = get_point(24)     # Right Hip
            right_knee = get_point(26)    # Right Knee
            right_ankle = get_point(28)   # Right Ankle
            
            # ✅ Compute gait-specific angles
            left_knee_angle = calculate_angle(left_hip, left_knee, left_ankle)
            right_knee_angle = calculate_angle(right_hip, right_knee, right_ankle)
            
            # ✅ Detect gait abnormalities
            risk_text = ""
            if left_knee_angle < 160 or right_knee_angle < 160:
                risk_text = "Abnormal Knee Angle: Risk of Unstable Gait!"
            elif abs(left_knee_angle - right_knee_angle) > 20:
                risk_text = "Asymmetric Knee Angles: Risk of Uneven Gait!"
            
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

# Route to stream webcam feed
@app.route("/webcam")
def webcam_feed():
    return Response(generate_frames(), mimetype='multipart/x-mixed-replace; boundary=frame')



if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
