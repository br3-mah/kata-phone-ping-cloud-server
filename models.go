package main

import "time"

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
