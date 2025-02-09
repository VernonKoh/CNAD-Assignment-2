# Fall Risk Self-Assessment Microservices Architecture

This project is designed using **Microservices Architecture** to support a **Fall Risk Self-Assessment System**. The system is split into distinct services to ensure scalability, maintainability, and independent deployment. The architecture supports fall risk assessment, cognitive training, chatbot interactions, and real-time video analysis.

## Table of Contents
1. [Design Considerations of the Microservices](#design-considerations)
2. [Architecture Diagram](#architecture-diagram)
3. [Instructions for Setting Up and Running Microservices](#setup-instructions)
4. [Why Microservices Architecture?](#why-microservices)
5. [System Requirements](#system-requirements)
6. [Known Issues and Limitations](#known-issues-and-limitations)

## Design Considerations of the Microservices
The platform was initially developed as a **local microservices-based architecture** and later scaled with **Docker** to enhance deployment efficiency. This transition showcases the evolution of our solution from a traditional microservices setup to a cloud-native, containerized infrastructure. Our approach focuses on gradual improvements while maintaining system reliability and accessibility.

### 1. Service Decomposition

The system is decomposed into distinct services that handle specific functionalities. Each service operates independently and interacts with others via **RESTful APIs**.

- **User Service**: Manages user registration, authentication, and profile management. Users can sign up, log in, and access their assessment history.
- **Assessment Service**: Handles the fall risk self-assessment questionnaire and scoring logic. Users answer a set of questions to evaluate their fall risk level.
- **Game Service**: Provides a cognitive training game (Memory Card Matching) to improve mental agility. Tracks game progress and performance.
- **LionBee Chatbot Service (DeepSeek API)**: A RESTful chatbot service that provides interactive support using speech-to-text technology. Supports four languages: English, Chinese, Malay, and Tamil.
- **Video Analysis Service & Posture Analysis**: Uses AI-driven motion tracking to analyze walking patterns and detect fall risks. Processes video data for movement stability assessment.

**Benefits of Service Decomposition**:
- **Scalability**: Each service can be scaled independently. For example, the **Video Analysis Service** can be scaled separately to handle higher computational demands.
- **Maintainability**: Services can be updated or maintained without affecting other parts of the system.

### 2. Inter-Service Communication

Services communicate using **RESTful APIs** with well-defined API contracts. This ensures that the system remains loosely coupled and flexible.

- **User Service API**: Manages user-related operations such as registration and authentication.
- **Assessment Service API**: Provides endpoints for managing fall risk questionnaires and results.
- **Game Service API**: Handles game-related logic and progress tracking.
- **Chatbot Service API**: Processes speech-to-text interactions and chatbot conversations.

By following RESTful principles, services can evolve independently without affecting the entire system.

### 3. Database Per Service

Each microservice maintains its own database to ensure **data isolation and independence**. This improves **fault tolerance** and prevents cross-service data dependencies.

- **User Service Database**: Stores user credentials, profiles, and history.
- **Assessment Service Database**: Stores questionnaire results and risk assessment data.
- **Game Service Database**: Maintains game progress and user scores.
- **Chatbot Service Database**: Manages chatbot interactions and user responses.
- **Video Analysis Service Database**: Stores motion analysis data and reports.

This ensures that if one database fails, the other services can continue operating without disruption.

### 4. Error Handling and Logging

Each service includes robust **error handling** and **logging mechanisms**:

- **Structured Logging**: Logs important events such as API requests, database operations, and errors.
- **Meaningful HTTP Responses**: Services return appropriate status codes (e.g., `400` for invalid input, `500` for server errors).
- **Service Resilience**: If one service fails, others can continue running independently.

## Architecture Diagram

Below is the architecture diagram illustrating the microservices and their interactions:

```plaintext
                           +------------------+
                           |   User Service   |
                           | (Authentication, |
                           | Profiles, History)|
                           +--------+---------+
                                    |
                          RESTful API |  
                                    |
          +-------------------------+------------------------+
          |                                                  |
+---------+------------+                          +---------+------------+
| Assessment Service  |                          |  Chatbot Service     |
| (Questionnaire,     |                          |  (Speech-to-Text)     |
| Risk Calculation)   |                          +----------------------+
+---------------------+                              
          |                                          |
       RESTful API                                  RESTful API
          |                                          |
+---------+------------+                          +---------+------------+
| Game Service        |                          |  Video Analysis       |
| (Memory Card Game)  |                          |  (Motion Tracking AI) |
+---------------------+                          +----------------------+
          |                                          |
       RESTful API                                  RESTful API
          |                                          |
+---------+------------+                          +---------+------------+
|  Database (Users)    |                          |  Database (Analysis)  |
+---------------------+                          +----------------------+

```

# Fall Risk Self-Assessment Microservices Architecture

## Key Takeaways:
- **Independent Microservices**: Each service runs independently with its own database.
- **Fault Isolation**: If one service fails, others continue operating without disruption.
- **Scalability**: Services can be individually scaled based on workload.

## Instructions for Setting Up and Running Microservices

### 1. Prerequisites

- Install **Go**, **Docker**, and **Python** (for video processing).
- Install **Pip** for Python dependencies:
  ```sh
  pip install opencv-python numpy mediapipe flask

## 2. Installation Steps

### Step 1: Start the Application

Run dependency management:
```sh
go mod tidy
```

## Instructions for Setting Up and Running Microservices

### 1. Run User Service and Game Service via Docker:

```sh
docker-compose build && docker-compose up
```

Start other services manually in separate terminal windows:

```sh
go run assessment-service/main.go
go run chat-service/main.go
python mediapipe_server.py
```

Step 2: Access the Application
Website: http://localhost:8081
API Endpoint: http://localhost:8081/api

## 5. System Requirements
To run the Fall Risk Self-Assessment Microservices Architecture system smoothly, make sure your environment meets the following requirements:

**Hardware:**
- Operating System: Windows, Linux, or macOS
- Memory: Minimum 4GB RAM (8GB recommended for optimal performance)
- Processor: Multi-core processor (preferably 2.0GHz or higher)
- Storage: Minimum 10GB free disk space

**Software:**
Docker: For containerizing and managing services.

**Go 1.16+:** For backend services, including the User, Assessment, and Game services.

**Python 3.7+:** For running the Video Analysis service.

**MediaPipe:** Python library for video processing.

Install using:

```sh
pip install opencv-python numpy mediapipe flask
```

## 6. Known Issues and Limitations
There are a few known issues and limitations in the system:

**1. FaceIO Public ID**
The FaceIO public ID needs to be renewed weekly as it is on a 7-day free trial. The current trial is valid until 16th February. Please contact Aaron (student developer) if you encounter any issues with the FaceIO service.

**2. LionBee Chatbot (DeepSeek API)**
The LionBee Chatbot utilizes the DeepSeek API. It uses an API key obtained from OpenRouter. Occasionally, the API may become disabled due to token limitations, which may cause the service to be unavailable at certain times. The API token needs to be manually enabled. If you encounter any errors, please contact Aaron so he can enable the API key.

**Why Microservices Architecture?**
1. Scalability: Individual services can be scaled based on demand. For example, Video Analysis Service can scale separately to handle increased processing loads.
2. Fault Isolation: If one service goes down, the others continue running. This ensures higher reliability and uptime.
3. Independent Databases: Each service manages its own database, improving data security and reducing dependencies.
4. Flexibility and Maintainability:
   - Independent Development: Teams can work on different services simultaneously.
   - Faster Deployments: Updates to one service do not require changes to others.
   - 
**Conclusion**
This microservices-based Fall Risk Self-Assessment System ensures a scalable, resilient, and modular solution. By leveraging RESTful APIs, independent databases, and containerized deployments, this architecture enhances flexibility, maintainability, and long-term system performance.
