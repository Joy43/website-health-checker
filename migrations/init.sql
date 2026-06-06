-- Create the database if it does not exist (fallback/safety check)
CREATE DATABASE IF NOT EXISTS health_checker;
USE health_checker;

-- Create health_checks table to verify schema setup and DB connectivity works
CREATE TABLE IF NOT EXISTS health_checks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert a test seed row
INSERT INTO health_checks (status) VALUES ('initialized');
