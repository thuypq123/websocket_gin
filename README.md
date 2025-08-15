# 🚀 Real-time WebSocket Chat & Comments System

A modern, real-time chat and commenting system built with Go, WebSocket, and SQLite. Features include multi-room chat, post commenting, and event-driven architecture.

## ✨ Features

- **🔄 Real-time Communication**: WebSocket-based instant messaging
- **💬 Multi-room Chat**: Support for multiple chat rooms (General, Tech, Random)
- **📝 Post Comments**: Real-time commenting system for posts
- **💾 Data Persistence**: All messages and comments saved to SQLite database
- **🎯 Event-driven Architecture**: Clean separation of concerns with event handlers
- **📡 RESTful API**: Complete REST API for posts, messages, and comments
- **🧪 Built-in Testing**: Test endpoints and WebSocket debugging tools
- **🎨 Modern UI**: Clean, responsive web interface

## 🏗️ Architecture

### Backend (Go)
- **Gin Gonic**: Web framework for HTTP routing
- **Gorilla WebSocket**: WebSocket implementation
- **SQLite**: Lightweight database for persistence
- **Event Router**: Centralized event handling system
- **Repository Pattern**: Clean data access layer

### Frontend (JavaScript)
- **Vanilla JS**: No framework dependencies
- **WebSocket Client**: Real-time communication
- **Responsive Design**: Works on desktop and mobile
- **Error Handling**: Comprehensive error reporting

## 📋 Prerequisites

- **Go 1.19+**: [Download Go](https://golang.org/dl/)
- **Git**: Version control
- **Modern Browser**: Chrome, Firefox, Safari, Edge

## 🚀 Quick Start

### 1. Clone Repository
```bash
git clone <your-repo-url>
cd ws
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Build Application
```bash
go build -o bin/server cmd/server/main.go
```

### 4. Run Server
```bash
./bin/server
```

### 5. Open Browser
- **Chat**: http://localhost:8080/chat
- **Posts**: http://localhost:8080/posts
- **API Health**: http://localhost:8080/api/v1/health

## 📚 API Documentation

### WebSocket Endpoint
```
ws://localhost:8080/ws?username=<your-username>
```

### WebSocket Events

#### Client → Server Events

**Join Room**
```json
{
  "type": "JOIN_ROOM",
  "room": "general",
  "user": "username"
}
```

**Send Chat Message**
```json
{
  "type": "CHAT_MESSAGE",
  "room": "general",
  "message": "Hello everyone!",
  "user": "username"
}
```

**Post Comment**
```json
{
  "type": "POST_COMMENT",
  "post_id": "post123",
  "comment": "Great post!",
  "user": "username"
}
```

#### Server → Client Events

**Room Joined Confirmation**
```json
{
  "type": "ROOM_JOINED",
  "room": "general",
  "user": "username"
}
```

**Chat Message Broadcast**
```json
{
  "type": "CHAT_MESSAGE",
  "room": "general",
  "message": "Hello everyone!",
  "user": "sender",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

**Comment Broadcast**
```json
{
  "type": "POST_COMMENT",
  "post_id": "post123",
  "comment": "Great post!",
  "user": "commenter",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

**Error Response**
```json
{
  "type": "ERROR",
  "message": "Error description"
}
```

### REST API Endpoints

#### Health Check
```http
GET /api/v1/health
```

#### Messages
```http
GET /api/v1/messages/{room}?limit=10    # Get recent messages
GET /api/v1/messages/recent             # Get all recent messages
```

#### Posts
```http
GET    /api/v1/posts                    # Get all posts
POST   /api/v1/posts                    # Create post
GET    /api/v1/posts/{id}               # Get specific post
PUT    /api/v1/posts/{id}               # Update post
DELETE /api/v1/posts/{id}               # Delete post
GET    /api/v1/posts/{id}/comments      # Get post comments
```

#### Testing Endpoints
```http
GET /api/v1/test/message?room=general&message=test&user=testuser
GET /api/v1/test/comment?post_id=1&comment=test&user=testuser
GET /api/v1/stats                       # WebSocket connection stats
```

## 🗂️ Project Structure

```
.
├── cmd/server/                 # Application entry point
│   └── main.go
├── internal/
│   ├── handlers/              # HTTP handlers
│   │   ├── enhanced_routes.go # Route definitions
│   │   ├── simple_chat.go     # Chat handlers
│   │   └── post_handler.go    # Post handlers
│   ├── repository/            # Data access layer
│   │   ├── message_repository.go
│   │   ├── comment_repository.go
│   │   └── post_repository.go
│   └── websocket/            # WebSocket implementation
│       ├── hub.go            # WebSocket hub (connection manager)
│       ├── client.go         # WebSocket client
│       ├── connection.go     # Connection handling
│       ├── event_router.go   # Event routing
│       └── handlers/         # Event handlers
│           ├── chat/         # Chat event handlers
│           ├── comments/     # Comment event handlers
│           ├── rooms/        # Room management
│           └── shared/       # Shared types/interfaces
├── templates/                # HTML templates
│   ├── chat.html
│   ├── post.html
│   ├── posts.html
│   └── index.html
├── static/                   # Static assets
│   ├── css/style.css
│   ├── js/simple-chat-fixed.js
│   └── test-websocket.html   # WebSocket testing tool
├── pkg/database/             # Database utilities
└── bin/                      # Compiled binaries
```

## 🧪 Testing

### WebSocket Testing Tool
Visit http://localhost:8080/static/test-websocket.html for a comprehensive WebSocket testing interface.

### Manual API Testing
```bash
# Test message sending
curl "http://localhost:8080/api/v1/test/message?room=general&message=Hello&user=testuser"

# Test comment posting
curl "http://localhost:8080/api/v1/test/comment?post_id=1&comment=Nice&user=testuser"

# Check connection stats
curl "http://localhost:8080/api/v1/stats"

# Get recent messages
curl "http://localhost:8080/api/v1/messages/general?limit=5"
```

## 🔧 Configuration

### Environment Variables
```bash
PORT=8080                    # Server port (default: 8080)
GIN_MODE=release            # Gin mode (debug/release)
```

### Database
- **Type**: SQLite
- **Location**: `./database.db`
- **Auto-migration**: Enabled
- **Sample data**: Auto-inserted on first run

## 🚀 Deployment

### Development
```bash
go run cmd/server/main.go
```

### Production Build
```bash
# Build for current OS
go build -o bin/server cmd/server/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o bin/server.exe cmd/server/main.go
```

### Docker (Optional)
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./server"]
```

## 🛠️ Development

### Adding New Event Types
1. **Define event in `events.go`**:
```go
const EventNewFeature = "NEW_FEATURE"
```

2. **Create handler in `handlers/`**:
```go
func (h *Handler) HandleNewFeature(client shared.ClientInterface, data []byte) error {
    // Implementation
}
```

3. **Register in event router**:
```go
case EventNewFeature:
    return r.newFeatureHandler.HandleNewFeature(client, messageBytes)
```

### Database Migrations
The application automatically handles database schema creation and updates on startup.

## 🐛 Troubleshooting

### Common Issues

**WebSocket Connection Failed**
- Check if server is running on correct port
- Verify firewall settings
- Check browser console for errors

**Messages Not Persisting**
- Check database file permissions
- Verify SQLite installation
- Check server logs for database errors

**Frontend Not Loading**
- Ensure static files are served correctly
- Check browser console for JavaScript errors
- Verify template files exist

### Debug Mode
Set `GIN_MODE=debug` for detailed logging:
```bash
export GIN_MODE=debug
./bin/server
```

### Logs Location
- **Console**: All logs printed to stdout
- **WebSocket Events**: Detailed event logging
- **Database Operations**: SQL query logging in debug mode

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-repo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-repo/discussions)
- **Documentation**: This README and inline code comments

## 🎯 Roadmap

- [ ] User authentication system
- [ ] Private messaging
- [ ] File upload support
- [ ] Message reactions/emojis
- [ ] Room moderation features
- [ ] Mobile app (React Native/Flutter)
- [ ] Redis integration for scaling
- [ ] Message encryption

---

**Built with ❤️ using Go, WebSocket, and modern web technologies**