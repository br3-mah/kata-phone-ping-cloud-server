package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
