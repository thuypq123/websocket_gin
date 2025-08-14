# WebSocket Chat Application

Một ứng dụng chat real-time được xây dựng bằng Go, Gin framework và WebSocket.

## Tính năng

- 🚀 Chat real-time với WebSocket
- 👥 Hỗ trợ nhiều phòng chat
- 🎨 Giao diện người dùng hiện đại và responsive
- 🔄 Tự động kết nối lại khi mất kết nối
- 📱 Tương thích với mobile
- ⚡ Hiệu suất cao và khả năng mở rộng

## Cấu trúc dự án

```
websocket/
├── cmd/server/           # Entry point của ứng dụng
│   └── main.go
├── internal/             # Internal packages
│   ├── handlers/         # HTTP handlers và routes
│   ├── models/          # Data structures
│   └── websocket/       # WebSocket hub và client logic
├── static/              # Static assets
│   ├── css/
│   └── js/
├── templates/           # HTML templates
├── go.mod
└── README.md
```

## Yêu cầu

- Go 1.21 hoặc cao hơn
- Các dependencies sẽ được tự động tải xuống khi chạy `go mod tidy`

## Cài đặt và chạy

1. Clone dự án và di chuyển vào thư mục:
   ```bash
   cd websocket
   ```

2. Tải dependencies:
   ```bash
   go mod tidy
   ```

3. Chạy server:
   ```bash
   go run cmd/server/main.go
   ```

4. Mở trình duyệt và truy cập:
   ```
   http://localhost:8080
   ```

## Sử dụng

1. Truy cập trang chủ tại `http://localhost:8080`
2. Nhập username và tên phòng chat (mặc định là "general")
3. Nhấn "Join Chat" để vào phòng chat
4. Bắt đầu chat với những người dùng khác trong cùng phòng!

## API Endpoints

- `GET /` - Trang chủ
- `GET /chat` - Giao diện chat
- `GET /ws` - WebSocket endpoint
- `GET /api/v1/health` - Health check

## Cấu hình

Ứng dụng sử dụng các biến môi trường sau:

- `PORT` - Port để chạy server (mặc định: 8080)

## Deployment

Để deploy ứng dụng:

```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Chạy
./bin/server
```

## Tính năng kỹ thuật

- **WebSocket Hub**: Quản lý tập trung các kết nối WebSocket
- **Room System**: Hỗ trợ nhiều phòng chat độc lập
- **Auto Reconnect**: Tự động kết nối lại khi mất kết nối
- **Message Types**: Hỗ trợ các loại message khác nhau (text, join, leave)
- **CORS Support**: Cấu hình CORS cho phát triển frontend riêng biệt

## Phát triển

Để phát triển thêm tính năng:

1. Thêm models mới trong `internal/models/`
2. Tạo handlers trong `internal/handlers/`
3. Cập nhật WebSocket logic trong `internal/websocket/`
4. Thêm frontend logic trong `static/js/`

## Đóng góp

1. Fork dự án
2. Tạo feature branch
3. Commit changes
4. Push to branch
5. Tạo Pull Request

## License

MIT License