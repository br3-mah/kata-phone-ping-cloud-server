package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type DeviceInfo struct {
	UUID      string            `json:"uuid" db:"uuid"`
	Hostname  string            `json:"hostname" db:"hostname"`
	OS        string            `json:"os" db:"os"`
	MAC       string            `json:"mac" db:"mac"`
	PublicIP  string            `json:"public_ip" db:"public_ip"`
	Geo       map[string]string `json:"geo"`
	Latitude  string            `json:"latitude" db:"latitude"`
	Longitude string            `json:"longitude" db:"longitude"`
}

type DeviceRecord struct {
	ID        int       `json:"id" db:"id"`
	UUID      string    `json:"uuid" db:"uuid"`
	Hostname  string    `json:"hostname" db:"hostname"`
	OS        string    `json:"os" db:"os"`
	MAC       string    `json:"mac" db:"mac"`
	PublicIP  string    `json:"public_ip" db:"public_ip"`
	Country   string    `json:"country" db:"country"`
	Region    string    `json:"region" db:"region"`
	City      string    `json:"city" db:"city"`
	Latitude  string    `json:"latitude" db:"latitude"`
	Longitude string    `json:"longitude" db:"longitude"`
	LastSeen  time.Time `json:"last_seen" db:"last_seen"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

var db *sql.DB

func main() {
	// Initialize database connection
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/ka_ping_db?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	// Create tables if they don't exist
	createTables()

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	api := r.Group("/api")
	{
		api.POST("/device-ping", handleDevicePing)
		api.GET("/devices", getDevices)
		api.GET("/device/:uuid", getDeviceByUUID)
		api.DELETE("/device/:uuid", deleteDevice)
	}

	// Web interface routes
	r.GET("/", serveIndex)
	r.Static("/static", "./static")

	log.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}

func createTables() {
	createTableSQL := `
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
		INDEX idx_last_seen (last_seen)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
	log.Println("Database tables created/verified successfully")
}

func handleDevicePing(c *gin.Context) {
	var deviceInfo DeviceInfo
	if err := c.ShouldBindJSON(&deviceInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract geo information for backward compatibility
	country := ""
	region := ""
	city := ""

	if deviceInfo.Geo != nil {
		country = deviceInfo.Geo["country"]
		region = deviceInfo.Geo["region"]
		city = deviceInfo.Geo["city"]
	}

	// Use latitude and longitude from the request
	latitude := deviceInfo.Latitude
	longitude := deviceInfo.Longitude

	// Check if device already exists
	var existingID int
	err := db.QueryRow("SELECT id FROM devices WHERE uuid = ?", deviceInfo.UUID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new device
		insertSQL := `
		INSERT INTO devices (uuid, hostname, os, mac, public_ip, country, region, city, latitude, longitude, last_seen)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`
		_, err = db.Exec(insertSQL, deviceInfo.UUID, deviceInfo.Hostname, deviceInfo.OS, deviceInfo.MAC,
			deviceInfo.PublicIP, country, region, city, latitude, longitude)
		if err != nil {
			log.Printf("Error inserting device: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert device"})
			return
		}
		log.Printf("New device registered: %s (%s) at %s,%s", deviceInfo.Hostname, deviceInfo.UUID, latitude, longitude)
	} else if err != nil {
		log.Printf("Error checking existing device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	} else {
		// Update existing device
		updateSQL := `
		UPDATE devices 
		SET hostname = ?, os = ?, mac = ?, public_ip = ?, country = ?, region = ?, city = ?, 
			latitude = ?, longitude = ?, last_seen = NOW(), updated_at = NOW()
		WHERE uuid = ?
		`
		_, err = db.Exec(updateSQL, deviceInfo.Hostname, deviceInfo.OS, deviceInfo.MAC, deviceInfo.PublicIP,
			country, region, city, latitude, longitude, deviceInfo.UUID)
		if err != nil {
			log.Printf("Error updating device: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
			return
		}
		log.Printf("Device updated: %s (%s) at %s,%s", deviceInfo.Hostname, deviceInfo.UUID, latitude, longitude)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Device ping received successfully",
		"timestamp": time.Now(),
	})
}

func getDevices(c *gin.Context) {
	rows, err := db.Query(`
		SELECT id, uuid, hostname, os, mac, public_ip, country, region, city, 
			   latitude, longitude, last_seen, created_at, updated_at 
		FROM devices 
		ORDER BY last_seen DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch devices"})
		return
	}
	defer rows.Close()

	var devices []DeviceRecord
	for rows.Next() {
		var device DeviceRecord
		err := rows.Scan(&device.ID, &device.UUID, &device.Hostname, &device.OS, &device.MAC,
			&device.PublicIP, &device.Country, &device.Region, &device.City,
			&device.Latitude, &device.Longitude, &device.LastSeen, &device.CreatedAt, &device.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan device"})
			return
		}
		devices = append(devices, device)
	}

	c.JSON(http.StatusOK, devices)
}

func getDeviceByUUID(c *gin.Context) {
	uuid := c.Param("uuid")
	var device DeviceRecord

	err := db.QueryRow(`
		SELECT id, uuid, hostname, os, mac, public_ip, country, region, city, 
			   latitude, longitude, last_seen, created_at, updated_at 
		FROM devices 
		WHERE uuid = ?
	`, uuid).Scan(&device.ID, &device.UUID, &device.Hostname, &device.OS, &device.MAC,
		&device.PublicIP, &device.Country, &device.Region, &device.City,
		&device.Latitude, &device.Longitude, &device.LastSeen, &device.CreatedAt, &device.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, device)
}

func deleteDevice(c *gin.Context) {
	uuid := c.Param("uuid")

	result, err := db.Exec("DELETE FROM devices WHERE uuid = ?", uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}

func serveIndex(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, indexHTML)
}

const indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ka-Ping Device Monitor</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: Arial, sans-serif; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; margin-bottom: 10px; }
        .stats { display: flex; justify-content: space-around; margin-bottom: 30px; }
        .stat-box { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; }
        .stat-number { font-size: 2em; font-weight: bold; color: #007bff; }
        .stat-label { color: #666; margin-top: 5px; }
        .devices-table { background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); overflow: hidden; }
        .table-header { background: #007bff; color: white; padding: 15px; }
        .table-content { overflow-x: auto; }
        table { width: 100%; border-collapse: collapse; }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #eee; }
        th { background: #f8f9fa; font-weight: bold; }
        .status-online { color: #28a745; }
        .status-offline { color: #dc3545; }
        .refresh-btn { background: #007bff; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; margin-bottom: 20px; }
        .refresh-btn:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Ka-Ping Device Monitor</h1>
            <p>Real-time device monitoring and tracking</p>
        </div>

        <div class="stats">
            <div class="stat-box">
                <div class="stat-number" id="totalDevices">0</div>
                <div class="stat-label">Total Devices</div>
            </div>
            <div class="stat-box">
                <div class="stat-number" id="onlineDevices">0</div>
                <div class="stat-label">Online (Last 10 min)</div>
            </div>
            <div class="stat-box">
                <div class="stat-number" id="offlineDevices">0</div>
                <div class="stat-label">Offline</div>
            </div>
        </div>

        <button class="refresh-btn" onclick="loadDevices()">Refresh</button>

        <div class="devices-table">
            <div class="table-header">
                <h3>Connected Devices</h3>
            </div>
            <div class="table-content">
                <table id="devicesTable">
                    <thead>
                        <tr>
                            <th>Hostname</th>
                            <th>UUID</th>
                            <th>OS</th>
                            <th>MAC Address</th>
                            <th>Public IP</th>
                            <th>Location</th>
                            <th>Last Seen</th>
                            <th>Status</th>
                        </tr>
                    </thead>
                    <tbody id="devicesTableBody">
                    </tbody>
                </table>
            </div>
        </div>
    </div>

    <script>
        async function loadDevices() {
            try {
                const response = await fetch('/api/devices');
                const devices = await response.json();
                
                const tbody = document.getElementById('devicesTableBody');
                tbody.innerHTML = '';

                let totalDevices = devices.length;
                let onlineDevices = 0;
                let offlineDevices = 0;

                devices.forEach(device => {
                    const row = document.createElement('tr');
                    const now = new Date();
                    const lastSeen = new Date(device.last_seen);
                    const timeDiff = (now - lastSeen) / (1000 * 60); // minutes
                    const isOnline = timeDiff <= 10;

                    if (isOnline) onlineDevices++;
                    else offlineDevices++;

                    const location = [device.city, device.region, device.country].filter(Boolean).join(', ') || 'Unknown';
                    const statusClass = isOnline ? 'status-online' : 'status-offline';
                    const statusText = isOnline ? 'Online' : 'Offline';

                    const coordinates = device.latitude && device.longitude ? 
                        device.latitude + ', ' + device.longitude : 'N/A';
                    const mapLink = device.latitude && device.longitude ? 
                        '<a href="https://maps.google.com/?q=' + device.latitude + ',' + device.longitude + '" target="_blank" style="color: #007bff; text-decoration: none;">View Map</a>' : 
                        'N/A';
                    
                    row.innerHTML = ` + "`" + `
                        <td>${device.hostname}</td>
                        <td>${device.uuid.substring(0, 8)}...</td>
                        <td>${device.os}</td>
                        <td>${device.mac}</td>
                        <td>${device.public_ip}</td>
                        <td>${location}<br><small>${coordinates}</small><br>${mapLink}</td>
                        <td>${lastSeen.toLocaleString()}</td>
                        <td class="${statusClass}">${statusText}</td>
                    ` + "`" + `;
                    tbody.appendChild(row);
                });

                document.getElementById('totalDevices').textContent = totalDevices;
                document.getElementById('onlineDevices').textContent = onlineDevices;
                document.getElementById('offlineDevices').textContent = offlineDevices;
            } catch (error) {
                console.error('Error loading devices:', error);
            }
        }

        // Load devices on page load
        loadDevices();

        // Auto-refresh every 30 seconds
        setInterval(loadDevices, 30000);
    </script>
</body>
</html>
`
