# README

## 1. Design Consideration

### 1.1 Overview

The **Fall Risk Self-Assessment Website** is designed as an interactive tool to help elderly individuals assess their risk of falling. Developed in collaboration with **Lion Befrienders**, this solution aims to provide a seamless and engaging experience that empowers seniors to take preventive actions against falls. The system integrates various assessment techniques, including questionnaires, video analysis, and cognitive games, to deliver a comprehensive evaluation.

The platform was initially developed as a **local microservices-based architecture** and later scaled with **Docker** to enhance deployment efficiency. This transition showcases the evolution of our solution from a traditional microservices setup to a cloud-native, containerized infrastructure. Our approach focuses on gradual improvements while maintaining system reliability and accessibility.

The key considerations in the design include:

- **User Accessibility**: Ensuring that the interface is easy to use for elderly users.
- **Modularity**: The system is divided into independent components for better maintainability.
- **Scalability**: The architecture allows for easy expansion and integration of new features.
- **Performance**: Efficient algorithms and database optimizations are implemented to improve response times.
- **Security**: Proper authentication and authorization mechanisms are in place to ensure data safety.

### 1.2 Key Features

- **Fall Risk Questionnaire**: A structured questionnaire that evaluates an elderly user's mobility, medical conditions, and daily activities to estimate their fall risk.
- **Video Analysis**: Utilizes AI-driven motion tracking to analyze walking patterns and detect early indicators of instability.
- **Memory Card Matching Game**: A gamified cognitive training tool aimed at improving memory retention and mental agility.
- **Chatbot (LionBee)**: A REST API-driven chatbot that provides interactive support using speech-to-text technology, available in English, Chinese, Malay, and Tamil.
- **Email Notifications**: Automatically alerts relevant personnel when an elderly user is identified as high-risk, ensuring timely intervention.
- **Scalable Deployment Evolution**: The system began as a fully local microservices architecture and has progressively incorporated Docker to enhance scalability, deployment, and maintainability.

### 1.3 Key Design Decisions

- **Phased Microservice Migration**: Initially developed with a fully local microservice architecture, the system has been gradually transitioned to Docker containers. This phased approach demonstrates improvements in scalability, deployment efficiency, and operational resilience, aligning with modern cloud-native best practices.
- **User-Friendly Interface**: The website follows a simple and intuitive UI/UX design for ease of navigation.
- **Multilingual Support**: Ensures inclusivity by supporting four languages for chatbot interactions.
- **RESTful API**: The backend is exposed as a REST API for easy integration with different frontend clients.
- **Database Selection**: MySQL was chosen for its reliability and scalability.
- **Deployment Consideration**: The system is containerized using Docker for easy deployment and portability.

## 2. Architecture Diagram

Below is the architecture diagram of the solution:

[Frontend (React/HTML)]  ---> [API Gateway]  ---> [Backend (Go/Python)]  ---> [Database]
                                               |                          
                                               v                          
                                         [Video Processing]               
                                               |                          
                                               v                          
                                       [AI Model for Analysis]        


*Explanation of components:*

- **Frontend**: User interface built with HTML/CSS/JS.
- **API Gateway**: Manages requests and routes them to appropriate backend services.
- **Backend**: Handles data processing, risk calculations, chatbot interactions, and game logic.
- **Database**: Stores risk assessment results, game scores, etc.
- **Video Processing**: Analyzes user movements for fall risk detection.
- **AI Model for Analysis**: Processes video analysis data.
- **Chat-Service**: A REST API handling chatbot interactions, running locally.
- **Game-Service**: Containerized microservice managing the memory card matching game.
- **User-Service**: Containerized microservice managing user authentication and profiles.

## 3. Setup & Running Instructions

### 3.1 Prerequisites

- Install necessary dependencies such as Docker, Node.js, Go, etc.
- Ensure the database is set up and running.

### 3.2 Installation Steps

4. Start the application:
   - Run dependency management:
     ```sh
     go mod tidy
     ```
   - Run **User Service** and **Game Service** via Docker:
     ```sh
     docker-compose build && docker-compose up
     ```
   - Start other services manually in separate terminal windows:
     ```sh
     go run assessment-service/main.go
     go run chat-service/main.go
     ```
5. Access the application:
   - Website: `http://localhost:8081`
   - API Endpoint: `http://localhost:8081/api`

---