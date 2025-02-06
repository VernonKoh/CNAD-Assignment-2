-- Step 1: Drop the user if exists (optional)
DROP USER IF EXISTS 'user'@'localhost';

-- Step 2: Create the user with a password
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';

-- Step 3: Grant all privileges on the 'elderly' database to 'user'
GRANT ALL PRIVILEGES ON elderly.* TO 'user'@'localhost';

-- Step 4: Flush privileges to apply changes
FLUSH PRIVILEGES;

-- Step 5: Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS elderly;

-- Step 6: Select the database
USE elderly;

-- Step 7: Create the users table if it doesn't exist
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'Basic',
    is_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    facial_id VARCHAR(255) NULL,  -- ✅ Added column for storing facial ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE users ADD COLUMN high_risk BOOLEAN DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS user_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    age INT,
    gender VARCHAR(10),
    address TEXT,
    phone_number VARCHAR(15),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-- Create table for storing doctor accounts separately
CREATE TABLE IF NOT EXISTS doctors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    license_number VARCHAR(50) UNIQUE NOT NULL,
    hospital VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert a sample doctor profile
INSERT INTO doctors (email, password, name, license_number, hospital, is_verified) 
VALUES ('doctor@example.com', 'doctor123', 'Dr. John Doe', 'DOC123456', 'General Hospital', TRUE);


CREATE TABLE IF NOT EXISTS game_scores (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    score INT NOT NULL,
    level INT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE game_scores ADD COLUMN time_taken INT NOT NULL DEFAULT 0;


-- Create the Assessments table
CREATE TABLE Assessments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Create the Questions table
CREATE TABLE Questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    assessment_id INT NOT NULL,
    question_text TEXT NOT NULL,
    type ENUM('mcq', 'text', 'number') NOT NULL,
    FOREIGN KEY (assessment_id) REFERENCES Assessments(id)
);

-- Create Options table
CREATE TABLE Options (
    id INT AUTO_INCREMENT PRIMARY KEY,
    assessment_id INT NOT NULL,
    question_id INT NOT NULL,
    option_text TEXT NOT NULL,
	risk_value INT NOT NULL DEFAULT 0,
    FOREIGN KEY (assessment_id) REFERENCES Assessments(id),
    FOREIGN KEY (question_id) REFERENCES Questions(id)
);

CREATE TABLE CompletedAssessments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    assessment_id INT NOT NULL,
    user_id INT NOT NULL,
    total_risk_score INT NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (assessment_id) REFERENCES Assessments(id)
);
CREATE TABLE SelectedOptions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    completed_id INT NOT NULL,
    option_id INT NOT NULL,
    FOREIGN KEY (completed_id) REFERENCES CompletedAssessments(id),
    FOREIGN KEY (option_id) REFERENCES Options(id)
);


-- Populate the Assessments table
INSERT INTO Assessments (name, description)
VALUES 
    ('Elderly Fall Risk Assessment', 'Assessment to determine fall risk for elderly individuals.');

-- Populate the Questions table
INSERT INTO Questions (assessment_id, question_text, type)
VALUES
    (1, 'What is your age range?', 'mcq'),
    (1, 'Do you experience dizziness when standing?', 'mcq'),
    (1, 'How often do you exercise?', 'mcq'),
    (1, 'Do you smoke?', 'mcq');

-- Populate the Options table
-- Insert options with risk values and assessment IDs
INSERT INTO Options (assessment_id, question_id, option_text, risk_value)
VALUES
    -- Options for question 1 (age range) in Assessment 1
    (1, 1, 'Under 50', 0),
    (1, 1, '50-60', 1),
    (1, 1, '60-70', 2),
    (1, 1, 'Above 70', 3),

    -- Options for question 2 (dizziness) in Assessment 1
    (1, 2, 'Yes', 3),
    (1, 2, 'No', 0),

    -- Options for question 3 (exercise frequency) in Assessment 1
    (1, 3, 'Never', 3),
    (1, 3, '1-2 times a week', 2),
    (1, 3, '3-4 times a week', 1),
    (1, 3, 'Daily', 0),

    -- Options for question 5 (smoking) in Assessment 1
    (1, 4, 'Yes', 3),
    (1, 4, 'No', 0);
    
    
-- Set Column to FALSE, once email is sent update to TRUE
ALTER TABLE CompletedAssessments ADD COLUMN notified BOOLEAN DEFAULT FALSE;