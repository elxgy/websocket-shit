# WebSocket Chat Server

A real-time chat server built with Go, WebSockets, and MongoDB. Supports up to 4 concurrent users with authentication and persistent message storage.

## Features

- **Real-time messaging** using WebSockets
- **User authentication** with MongoDB
- **Persistent message storage** 
- **Support for up to 4 concurrent users**
- **CORS enabled** for frontend integration
- **Docker containerized** for easy deployment
- **Default users** for testing

## API Endpoints

### Authentication
- `POST /login` - User login
- `GET /health` - Health check

### WebSocket
- `GET /ws?username=<username>` - WebSocket connection (requires authentication first)

## Quick Start

### Prerequisites
- Go 1.21+
- MongoDB instance (Atlas or local)
- Docker (for deployment)

### Environment Setup

1. Create a `.env` file with your MongoDB URI:
```bash
MONGODB_URI="your_mongodb_connection_string"
PORT=8080
```

### Local Development

1. **Install dependencies:**
```bash
go mod download
```

2. **Create default users:**
```bash
cd scripts
go run setup_users.go ../models.go ../database.go
```

3. **Run the server:**
```bash
go run .
```

The server will start on port 8080 with the following endpoints:
- Health check: `http://localhost:8080/health`
- Login: `http://localhost:8080/login`
- WebSocket: `ws://localhost:8080/ws?username=<username>`

### Default Test Users

The system comes with 4 default users for testing:

| Username | Password    |
|----------|-------------|
| alice    | password123 |
| bob      | password123 |
| charlie  | password123 |
| diana    | password123 |

### Docker Deployment

1. **Build the Docker image:**
```bash
docker build -t chat-server .
```

2. **Run with environment variables:**
```bash
docker run -p 8080:8080 -e MONGODB_URI="your_mongodb_uri" chat-server
```

### Railway Deployment

1. **Push to GitHub repository**
2. **Connect Railway to your repository**
3. **Set environment variable:**
   - `MONGODB_URI`: Your MongoDB connection string
4. **Deploy**

The Railway deployment will automatically use the Dockerfile for containerized deployment.

## API Usage Examples

### Login Request
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "password123"}'
```

### WebSocket Connection (JavaScript)
```javascript
// First login to verify credentials
const loginResponse = await fetch('http://localhost:8080/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'alice', password: 'password123' })
});

if (loginResponse.ok) {
  // Connect to WebSocket
  const ws = new WebSocket('ws://localhost:8080/ws?username=alice');
  
  ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('Received:', message);
  };
  
  // Send a message
  ws.send(JSON.stringify({
    type: 'message',
    content: 'Hello, world!'
  }));
}
```

## Message Format

### Incoming Messages (from client)
```json
{
  "type": "message",
  "content": "Hello, world!"
}
```

### Outgoing Messages (to clients)
```json
{
  "type": "message",
  "username": "alice",
  "content": "Hello, world!",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### System Messages
```json
{
  "type": "user_joined",
  "username": "alice",
  "content": "alice joined the chat",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Database Schema

### Users Collection
```json
{
  "_id": "ObjectId",
  "username": "alice",
  "password": "hashed_password"
}
```

### Messages Collection
```json
{
  "_id": "ObjectId", 
  "username": "alice",
  "content": "Hello, world!",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Frontend Integration

This server is designed to work with a separate frontend application. Key integration points:

- **CORS enabled** for cross-origin requests
- **REST API** for authentication
- **WebSocket** for real-time messaging
- **JSON message format** for easy parsing

## Limitations

- **Maximum 4 concurrent users**
- **Simple username/password authentication** (no JWT tokens)
- **No user registration** (uses predefined users)
- **No private messaging** (broadcast only)

## Development

### Project Structure
```
├── main.go           # Server setup and HTTP handlers
├── hub.go           # WebSocket connection manager  
├── client.go        # WebSocket client handler
├── database.go      # MongoDB connection and operations
├── models.go        # Data structures
├── scripts/         # Utility scripts
│   └── setup_users.go # Create default users
├── Dockerfile       # Container configuration
├── go.mod          # Go module dependencies
└── .env            # Environment variables
```

### Adding New Features

1. **User Registration**: Extend the `/register` endpoint in `main.go`
2. **Private Messaging**: Modify `hub.go` to support targeted messages
3. **Message History**: Add pagination to `GetRecentMessages` in `database.go`
4. **User Roles**: Extend the `User` model in `models.go`

## Troubleshooting

### Common Issues

1. **MongoDB Connection**: Verify your `MONGODB_URI` is correct
2. **CORS Errors**: Ensure the server is running with CORS enabled
3. **WebSocket Connection**: Check that authentication was successful first
4. **Port Conflicts**: Change the `PORT` environment variable if needed

### Debug Mode

Enable verbose logging by adding debug statements:
```go
log.Printf("Debug: %+v", variableName)
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
websocket chat server for study and practice using go, docker and railway deploy
