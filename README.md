# Ka-Ping Server

A Go backend server with MySQL database integration to receive device information from ka-ping clients.

## Prerequisites

- Go 1.20 or later
- XAMPP with MySQL running
- MySQL accessible on localhost:3306

## Setup Instructions

### 1. Database Setup

1. Start XAMPP and ensure MySQL is running
2. Open phpMyAdmin or MySQL command line
3. Run the SQL script to create the database:
   ```bash
   mysql -u root -p < setup_database.sql
   ```
   Or execute the SQL commands in `setup_database.sql` manually through phpMyAdmin

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Build and Run

```bash
# Build the server
go build -o ka-ping-server.exe

# Run the server
./ka-ping-server.exe
```

The server will start on port 8080.

## API Endpoints

### Device Management
- `POST /api/device-ping` - Receive device ping data (used by ka-ping clients)
- `GET /api/devices` - Get all devices
- `GET /api/device/:uuid` - Get specific device by UUID
- `DELETE /api/device/:uuid` - Delete device by UUID

### Web Interface
- `GET /` - Web dashboard showing device status and information

## Web Dashboard

Visit `http://localhost:8080` to access the web dashboard that shows:
- Total devices count
- Online devices (last seen within 10 minutes)
- Offline devices
- Device details table with real-time updates

## Database Schema

The `devices` table contains:
- `id` - Auto-increment primary key
- `uuid` - Unique device identifier
- `hostname` - Device hostname
- `os` - Operating system information
- `mac` - MAC address
- `public_ip` - Public IP address
- `country`, `region`, `city` - Location information
- `latitude`, `longitude` - GPS coordinates
- `last_seen` - Last ping timestamp
- `created_at`, `updated_at` - Record timestamps

## Usage with Ka-Ping Client

To connect your ka-ping client to this server:

1. Update the endpoint URL in the ka-ping client configuration:
   ```go
   Endpoint: "http://localhost:8080/api/device-ping"
   ```

2. The client will automatically start sending device information every 5 minutes

## Features

- **Real-time monitoring**: Device status updates in real-time
- **Geolocation tracking**: Stores IP-based location data
- **Device persistence**: Maintains device records and updates
- **Web interface**: Easy-to-use dashboard for monitoring
- **RESTful API**: Full CRUD operations for device management
- **Auto-refresh**: Dashboard updates every 30 seconds

## Configuration

The server uses these default settings:
- **Port**: 8080
- **Database**: ka_ping_db
- **MySQL**: localhost:3306 (root user, no password)

To modify these settings, edit the connection string in `main.go`:
```go
db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/ka_ping_db?charset=utf8mb4&parseTime=True&loc=Local")
```

## Troubleshooting

1. **Database Connection Issues**:
   - Ensure MySQL is running in XAMPP
   - Check if the database `ka_ping_db` exists
   - Verify MySQL credentials

2. **Port Already in Use**:
   - Change the port in `main.go`: `r.Run(":8081")`

3. **CORS Issues**:
   - The server includes CORS headers for cross-origin requests
   - Modify the CORS middleware if needed

## Development

To extend functionality:
1. Add new API endpoints in `main.go`
2. Modify the database schema in `setup_database.sql`
3. Update the web interface HTML/CSS/JavaScript as needed

## Security Notes

- This is a development setup with no authentication
- For production use, implement proper authentication and authorization
- Use environment variables for sensitive configuration
- Enable HTTPS for production deployments
