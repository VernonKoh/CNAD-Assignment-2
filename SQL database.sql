CREATE USER 'user'@'localhost' IDENTIFIED BY
'password';
GRANT ALL ON *.* TO 'user'@'localhost'

-- Step 1: Create the database
CREATE DATABASE IF NOT EXISTS car_sharing;

-- Step 2: Use the database
USE car_sharing;
select * from users;

-- Step 3: Create the courses table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'Basic',
    is_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS membership_tiers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    hourly_rate_discount DECIMAL(5, 2) DEFAULT 0.0,
    priority_access BOOLEAN DEFAULT FALSE,
    booking_limit INT DEFAULT 0
);
