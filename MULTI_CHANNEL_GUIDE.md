# 🚀 Multi-Channel WebSocket System - Implementation Guide

## 📋 Tổng quan

Hệ thống Multi-Channel WebSocket cho phép một client kết nối WebSocket duy nhất có thể subscribe và nhận messages từ nhiều kênh khác nhau cùng lúc.

## 🏗️ Kiến trúc

### **Unified Channel System**
```
Client → Single WebSocket → Hub → Channel Manager → Multiple Channels
                                                   ├── chat:general
                                                   ├── chat:tech  
                                                   ├── post:123
                                                   ├── user:john
                                                   └── system:alerts
```

### **Channel Naming Convention**
```
Format: {type}:{identifier}

Examples:
- chat:general      (Chat room general)
- chat:tech         (Chat room tech) 
- post:123          (Comments for post 123)
- user:john         (Personal notifications for user john)
- system:alerts     (System-wide alerts)
- private:user1-user2 (Private conversation)
```

## 🎯 Các tính năng mới

### **1. Multi-Channel Subscription**
- Một WebSocket connection có thể subscribe nhiều channels
- Dynamic subscribe/unsubscribe channels
- Permission-based access control

### **2. Unified Event System**
- Tất cả events được route thông qua EventRouter
- Support cả legacy và new event types
- Backward compatibility với existing code

### **3. Channel Types**
- **chat**: Public chat rooms
- **post**: Post comment subscriptions  
- **user**: Private user notifications
- **system**: System-wide alerts
- **private**: Private conversations

## 📡 API Events

### **Client → Server Events**

#### **SUBSCRIBE**
```json
{
  "type": "SUBSCRIBE",
  "channels": ["chat:general", "post:123", "user:john"]
}
```

#### **UNSUBSCRIBE** 
```json
{
  "type": "UNSUBSCRIBE",
  "channels": ["chat:tech"]
}
```

#### **LIST_SUBSCRIPTIONS**
```json
{
  "type": "LIST_SUBSCRIPTIONS"
}
```

### **Server → Client Events**

#### **SUBSCRIPTION_RESPONSE**
```json
{
  "type": "SUBSCRIPTION_RESPONSE",
  "action": "subscribed",
  "channels": ["chat:general", "post:123"],
  "success": true
}
```

## 🔌 Frontend Usage

### **WebSocket Connection**
```javascript
const wsAdapter = new WebSocketAdapter({
    url: 'ws://localhost:8080/ws',
    username: 'john_doe'
});

wsAdapter.connect();
```

### **Subscribe to Multiple Channels**
```javascript
// Subscribe to multiple channels at once
wsAdapter.subscribe([
    'chat:general',      // General chat room
    'chat:tech',         // Tech chat room  
    'post:123',          // Comments for post 123
    'user:john_doe',     // Personal notifications
    'system:alerts'      // System alerts
]);
```

### **Handle Events from Different Channels**
```javascript
// Handle chat messages
wsAdapter.on('CHAT_MESSAGE', (data) => {
    console.log(`Chat in ${data.room}: ${data.message}`);
});

// Handle post comments
wsAdapter.on('POST_COMMENT', (data) => {
    console.log(`Comment on post ${data.post_id}: ${data.comment}`);
});

// Handle subscription responses
wsAdapter.on('SUBSCRIPTION_RESPONSE', (data) => {
    if (data.action === 'subscribed') {
        console.log('Subscribed to:', data.channels);
    }
});
```

### **Send Messages to Specific Channels**
```javascript
// Send chat message
wsAdapter.sendEvent('CHAT_MESSAGE', 'send', {
    type: 'CHAT_MESSAGE',
    room: 'general',
    message: 'Hello everyone!'
});

// Send post comment
wsAdapter.sendEvent('POST_COMMENT', 'send', {
    type: 'POST_COMMENT',
    post_id: '123',
    comment: 'Great post!'
});
```

## 🌐 Demo Pages

### **1. Multi-Channel Demo**
```
http://localhost:8080/demo
```
- Interactive demo cho multi-channel subscription
- Real-time message testing
- Channel management UI
- WebSocket statistics

### **2. Legacy Chat Demo**
```
http://localhost:8080/chat
```
- Original chat demo (still works)
- Backward compatibility testing

### **3. Posts Demo**
```
http://localhost:8080/posts
```
- Post commenting system
- Real-time comment updates

## 🔧 Backend Implementation

### **1. Channel Manager**
```go
// Subscribe to channel
hub.SubscribeToChannel(client, "chat:general")

// Broadcast to channel
hub.BroadcastToChannel("chat:general", event)

// Get client subscriptions
subscriptions := hub.GetClientSubscriptions(client)
```

### **2. Event Routing**
```go
// Events are automatically routed based on type
switch eventType {
case "SUBSCRIBE":
    return subscriptionHandler.HandleSubscribe(client, data)
case "CHAT_MESSAGE":
    return chatHandler.HandleChatMessage(client, data)
case "POST_COMMENT":
    return commentHandler.HandlePostComment(client, data)
}
```

### **3. Permission Control**
```go
// Validate channel access
func validateChannelAccess(client ClientInterface, channel *Channel) error {
    switch channel.Type {
    case ChannelTypeUser:
        // User channels are private to the user
        if channel.Identifier != client.GetUsername() {
            return fmt.Errorf("access denied")
        }
    case ChannelTypePrivate:
        // Check if user is participant
        if !canAccessPrivateChannel(client, channel.Identifier) {
            return fmt.Errorf("access denied")
        }
    }
    return nil
}
```

## 📊 Monitoring & Statistics

### **WebSocket Stats API**
```bash
curl http://localhost:8080/api/v1/ws/stats
```

Response:
```json
{
    "websocket_stats": {
        "total_clients": 5,
        "chat_rooms": 3,
        "post_subscribers": 2
    },
    "channel_stats": {
        "chat:general": 3,
        "chat:tech": 2,
        "post:123": 1,
        "user:john": 1
    },
    "supported_events": [
        "JOIN_ROOM", "CHAT_MESSAGE", "POST_COMMENT",
        "SUBSCRIBE", "UNSUBSCRIBE", "LIST_SUBSCRIPTIONS"
    ]
}
```

## 🔄 Migration từ Legacy System

### **Backward Compatibility**
- Existing code vẫn hoạt động bình thường
- Legacy events (`JOIN_ROOM`, `CHAT_MESSAGE`, `POST_COMMENT`) được support
- Automatic fallback to legacy methods nếu channel system fails

### **Migration Strategy**
1. **Phase 1**: Deploy new system với backward compatibility
2. **Phase 2**: Update frontend để sử dụng new subscription events
3. **Phase 3**: Gradually migrate existing features
4. **Phase 4**: Remove legacy code (optional)

## 🧪 Testing

### **Manual Testing**
1. Open demo page: `http://localhost:8080/demo`
2. Enter username và connect
3. Subscribe to multiple channels
4. Send messages to different channels
5. Open multiple browser tabs để test real-time sync

### **API Testing**
```bash
# Test subscription via curl (requires WebSocket client)
wscat -c ws://localhost:8080/ws?username=testuser

# Send subscription event
{"type":"SUBSCRIBE","channels":["chat:general","post:123"]}

# Send chat message  
{"type":"CHAT_MESSAGE","room":"general","message":"Hello"}
```

## 🚀 Production Considerations

### **Performance**
- Channel subscriptions are stored in memory maps
- Thread-safe với RWMutex
- Efficient broadcasting với buffered channels

### **Scalability**
- Horizontal scaling với Redis pub/sub (future enhancement)
- Connection pooling và load balancing
- Channel-based sharding strategies

### **Security**
- Permission validation cho private channels
- Input validation cho channel names
- Rate limiting cho subscription requests

## 📈 Future Enhancements

- [ ] Redis integration cho distributed scaling
- [ ] Message persistence cho offline users  
- [ ] Channel history và replay functionality
- [ ] Advanced permission system với roles
- [ ] WebRTC integration cho voice/video channels
- [ ] Mobile app support
- [ ] Message encryption cho private channels

## 🎯 Kết luận

Multi-Channel WebSocket System cung cấp:

✅ **Flexible Subscription**: Subscribe nhiều channels với một connection
✅ **Real-time Communication**: Instant messaging across different domains  
✅ **Scalable Architecture**: Clean separation of concerns
✅ **Backward Compatibility**: Không breaking existing functionality
✅ **Production Ready**: Thread-safe, error handling, monitoring

Hệ thống này cho phép xây dựng các ứng dụng real-time phức tạp với architecture sạch và dễ maintain! 🎉
