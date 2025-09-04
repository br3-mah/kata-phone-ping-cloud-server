-- Create database
CREATE DATABASE IF NOT EXISTS ka_ping_db;
USE ka_ping_db;

-- Create devices table
CREATE TABLE IF NOT EXISTS devices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    hostname VARCHAR(255) NOT NULL,
    os VARCHAR(100) NOT NULL,
    mac VARCHAR(17) NOT NULL,
    public_ip VARCHAR(45) NOT NULL,
    country VARCHAR(100),
    region VARCHAR(100),
    city VARCHAR(100),
    latitude VARCHAR(20),
    longitude VARCHAR(20),
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uuid (uuid),
    INDEX idx_last_seen (last_seen),
    INDEX idx_hostname (hostname)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create a view for online devices (last seen within 10 minutes)
CREATE OR REPLACE VIEW online_devices AS
SELECT *
FROM devices
WHERE last_seen >= DATE_SUB(NOW(), INTERVAL 10 MINUTE);

-- Create a view for offline devices (last seen more than 10 minutes ago)
CREATE OR REPLACE VIEW offline_devices AS
SELECT *
FROM devices
WHERE last_seen < DATE_SUB(NOW(), INTERVAL 10 MINUTE);

-- Sample queries you can use:
-- SELECT COUNT(*) as total_devices FROM devices;
-- SELECT COUNT(*) as online_devices FROM online_devices;
-- SELECT COUNT(*) as offline_devices FROM offline_devices;
-- SELECT * FROM devices ORDER BY last_seen DESC LIMIT 10;
