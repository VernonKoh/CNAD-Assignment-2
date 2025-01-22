CREATE USER 'user'@'localhost' IDENTIFIED BY
'password';
GRANT ALL ON *.* TO 'user'@'localhost'

-- Step 1: Create the database
CREATE DATABASE IF NOT EXISTS elderly;

-- Step 2: Use the database
USE elderly;
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


