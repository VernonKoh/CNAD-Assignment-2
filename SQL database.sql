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


CREATE TABLE IF NOT EXISTS assessments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    assessment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    score INT,
    risk_level VARCHAR(50),
    notes TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
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


