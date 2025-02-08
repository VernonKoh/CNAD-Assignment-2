-- Step 1: Drop the user if it exists
DROP USER IF EXISTS 'user'@'localhost';

-- Step 2: Create the user with a password (only if it doesn't exist)
CREATE USER IF NOT EXISTS 'user'@'localhost' IDENTIFIED BY 'password';

-- Step 3: Grant privileges to the user
GRANT ALL PRIVILEGES ON elderly.* TO 'user'@'localhost';

-- Step 4: Flush privileges
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
    facial_id VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Check if the column 'high_risk' exists before adding it
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
                   WHERE TABLE_NAME = 'users' 
                   AND COLUMN_NAME = 'high_risk' 
                   AND TABLE_SCHEMA = DATABASE());

-- Add column only if it does not exist
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN high_risk BOOLEAN DEFAULT FALSE;', 'SELECT "Column already exists";');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Step 9: Create user_details table if it doesn't exist
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

-- Step 10: Create doctors table if it doesn't exist
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

-- Step 11: Insert a sample doctor profile only if it doesn't already exist
INSERT INTO doctors (email, password, name, license_number, hospital, is_verified) 
SELECT 'doctor@example.com', 'doctor123', 'Dr. John Doe', 'DOC123456', 'General Hospital', TRUE
WHERE NOT EXISTS (SELECT 1 FROM doctors WHERE email = 'doctor@example.com');

-- Step 12: Create game_scores table if it doesn't exist
CREATE TABLE IF NOT EXISTS game_scores (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    score INT NOT NULL,
    level INT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Step 13: Add column only if it doesn't exist
-- Check if the column 'time_taken' exists before adding it
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
                   WHERE TABLE_NAME = 'game_scores' 
                   AND COLUMN_NAME = 'time_taken' 
                   AND TABLE_SCHEMA = DATABASE());

-- Add column only if it does not exist
SET @query = IF(@col_exists = 0, 'ALTER TABLE game_scores ADD COLUMN time_taken INT NOT NULL DEFAULT 0;', 'SELECT "Column already exists";');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Step 14: Create the Assessments table
CREATE TABLE IF NOT EXISTS Assessments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Step 15: Create the Questions table
CREATE TABLE IF NOT EXISTS Questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    assessment_id INT NOT NULL,
    question_text TEXT NOT NULL,
    type ENUM('mcq', 'text', 'number') NOT NULL,
    FOREIGN KEY (assessment_id) REFERENCES Assessments(id)
);

-- Step 16: Create Options table
CREATE TABLE IF NOT EXISTS Options (
    id INT AUTO_INCREMENT PRIMARY KEY,
    assessment_id INT NOT NULL,
    question_id INT NOT NULL,
    option_text TEXT NOT NULL,
    risk_value INT NOT NULL DEFAULT 0,
    FOREIGN KEY (assessment_id) REFERENCES Assessments(id),
    FOREIGN KEY (question_id) REFERENCES Questions(id)
);

-- Step 17: Create CompletedAssessments table
CREATE TABLE IF NOT EXISTS CompletedAssessments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    assessment_id INT NOT NULL,
    user_id INT NOT NULL,
    total_risk_score INT NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (assessment_id) REFERENCES Assessments(id)
);

-- Step 18: Create SelectedOptions table
CREATE TABLE IF NOT EXISTS SelectedOptions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    completed_id INT NOT NULL,
    option_id INT NOT NULL,
    FOREIGN KEY (completed_id) REFERENCES CompletedAssessments(id),
    FOREIGN KEY (option_id) REFERENCES Options(id)
);

-- Step 19: Insert assessment only if it doesn’t exist
INSERT INTO Assessments (id, name, description)
VALUES (1, 'Elderly Fall Risk Assessment', 'Assessment to determine fall risk for elderly individuals.')
ON DUPLICATE KEY UPDATE name = name;

-- Step 20: Insert questions only if they don’t exist
INSERT INTO Questions (id, assessment_id, question_text, type)
VALUES
    (1, 1, 'What is your age range?', 'mcq'),
    (2, 1, 'Do you experience dizziness when standing?', 'mcq'),
    (3, 1, 'How often do you exercise?', 'mcq'),
    (4, 1, 'Do you smoke?', 'mcq')
ON DUPLICATE KEY UPDATE question_text = question_text;

-- Step 21: Insert options only if they don’t exist
INSERT INTO Options (id, assessment_id, question_id, option_text, risk_value)
VALUES
    (1, 1, 1, 'Under 50', 0),
    (2, 1, 1, '50-60', 1),
    (3, 1, 1, '60-70', 2),
    (4, 1, 1, 'Above 70', 3),

    (5, 1, 2, 'Yes', 3),
    (6, 1, 2, 'No', 0),

    (7, 1, 3, 'Never', 3),
    (8, 1, 3, '1-2 times a week', 2),
    (9, 1, 3, '3-4 times a week', 1),
    (10, 1, 3, 'Daily', 0),

    (11, 1, 4, 'Yes', 3),
    (12, 1, 4, 'No', 0)
ON DUPLICATE KEY UPDATE option_text = option_text;

-- Step 22: Add column only if it doesn't exist
-- Set Column to FALSE, once email is sent update to TRUE
-- Check if the column 'notified' exists before adding it
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
                   WHERE TABLE_NAME = 'CompletedAssessments' 
                   AND COLUMN_NAME = 'notified' 
                   AND TABLE_SCHEMA = DATABASE());

-- Add column only if it does not exist
SET @query = IF(@col_exists = 0, 'ALTER TABLE CompletedAssessments ADD COLUMN notified BOOLEAN DEFAULT FALSE;', 'SELECT "Column already exists";');
PREPARE stmt FROM @query;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
