class SimpleChatApp {
    constructor() {
        this.ws = null;
        this.username = '';
        this.currentRoom = '';
        this.isConnected = false;
        
        this.initializeElements();
        this.bindEvents();
        this.requestUsername();
    }

    initializeElements() {
        this.usernameModal = document.getElementById('usernameModal');
        this.usernameInput = document.getElementById('usernameInput');
        this.usernameSubmit = document.getElementById('usernameSubmit');
        this.roomSelect = document.getElementById('roomSelect');
        this.messagesDiv = document.getElementById('messages');
        this.messageInput = document.getElementById('messageInput');
        this.sendButton = document.getElementById('sendButton');
        this.connectionStatus = document.getElementById('connectionStatus');
    }

    bindEvents() {
        this.usernameSubmit.addEventListener('click', () => this.setUsername());
        this.usernameInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.setUsername();
        });
        
        this.roomSelect.addEventListener('change', (e) => this.joinRoom(e.target.value));
        this.sendButton.addEventListener('click', () => this.sendMessage());
        this.messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.sendMessage();
        });
    }

    requestUsername() {
        this.usernameModal.style.display = 'flex';
        this.usernameInput.focus();
    }

    setUsername() {
        const username = this.usernameInput.value.trim();
        if (!username) {
            alert('Please enter a username');
            return;
        }
        
        this.username = username;
        this.usernameModal.style.display = 'none';
        this.connect();
    }

    connect() {
        this.updateConnectionStatus('connecting');
        
        const wsUrl = `ws://localhost:8080/ws?username=${encodeURIComponent(this.username)}`;
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = () => {
            this.isConnected = true;
            this.updateConnectionStatus('connected');
            console.log('✅ Connected to WebSocket');
            
            // Auto-join general room
            this.joinRoom('general');
        };
        
        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleEvent(data);
            } catch (error) {
                console.error('❌ Error parsing message:', error);
            }
        };
        
        this.ws.onclose = () => {
            this.isConnected = false;
            this.updateConnectionStatus('disconnected');
            console.log('❌ WebSocket disconnected');
            
            // Try to reconnect after 3 seconds
            setTimeout(() => this.connect(), 3000);
        };
        
        this.ws.onerror = (error) => {
            console.error('❌ WebSocket error:', error);
            this.updateConnectionStatus('error');
        };
    }

    joinRoom(roomName) {
        if (!this.isConnected || !roomName) return;
        
        this.currentRoom = roomName;
        this.roomSelect.value = roomName;
        
        // Clear messages when switching rooms
        this.messagesDiv.innerHTML = '';
        
        // Send join room event
        this.sendEvent({
            type: 'JOIN_ROOM',
            room: roomName,
            user: this.username
        });
        
        console.log(`🏠 Joining room: ${roomName}`);
    }

    sendMessage() {
        const message = this.messageInput.value.trim();
        if (!message || !this.isConnected || !this.currentRoom) {
            return;
        }
        
        // Send chat message event
        this.sendEvent({
            type: 'CHAT_MESSAGE',
            room: this.currentRoom,
            message: message,
            user: this.username
        });
        
        this.messageInput.value = '';
    }

    sendEvent(event) {
        if (!this.isConnected) {
            console.warn('⚠️ Not connected, cannot send event');
            return false;
        }
        
        try {
            this.ws.send(JSON.stringify(event));
            return true;
        } catch (error) {
            console.error('❌ Error sending event:', error);
            return false;
        }
    }

    handleEvent(event) {
        console.log('📨 Received event:', event);
        
        switch (event.type) {
            case 'ROOM_JOINED':
                this.handleRoomJoined(event);
                break;
            case 'CHAT_MESSAGE':
                this.handleChatMessage(event);
                break;
            case 'POST_COMMENT':
                this.handlePostComment(event);
                break;
            case 'ERROR':
                this.handleError(event);
                break;
            default:
                console.warn('⚠️ Unknown event type:', event.type);
        }
    }

    handleRoomJoined(event) {
        this.addSystemMessage(`✅ Joined room: ${event.room} as ${event.user}`);
    }

    handleChatMessage(event) {
        // Only show messages for current room
        if (event.room === this.currentRoom) {
            this.addChatMessage(event.user, event.message);
        }
    }

    handlePostComment(event) {
        // Handle post comments (could be displayed in a separate section)
        this.addSystemMessage(`💬 ${event.user} commented on post ${event.post_id}: ${event.comment}`);
    }

    handleError(event) {
        console.error('❌ Server error:', event.message);
        this.addSystemMessage(`❌ Error: ${event.message}`);
    }

    addChatMessage(username, message) {
        const messageElement = document.createElement('div');
        messageElement.className = 'message';
        
        const isOwnMessage = username === this.username;
        if (isOwnMessage) {
            messageElement.classList.add('own-message');
        }
        
        const timestamp = new Date().toLocaleTimeString();
        messageElement.innerHTML = `
            <div class="message-header">
                <span class="username">${this.escapeHtml(username)}</span>
                <span class="timestamp">${timestamp}</span>
            </div>
            <div class="message-content">${this.escapeHtml(message)}</div>
        `;
        
        this.messagesDiv.appendChild(messageElement);
        this.scrollToBottom();
    }

    addSystemMessage(message) {
        const messageElement = document.createElement('div');
        messageElement.className = 'message system-message';
        messageElement.textContent = message;
        
        this.messagesDiv.appendChild(messageElement);
        this.scrollToBottom();
    }

    updateConnectionStatus(status) {
        const statusElement = this.connectionStatus;
        statusElement.className = `connection-status ${status}`;
        
        const statusText = {
            'connecting': '🔄 Connecting...',
            'connected': '✅ Connected',
            'disconnected': '❌ Disconnected',
            'error': '⚠️ Error'
        };
        
        statusElement.textContent = statusText[status] || status;
    }

    scrollToBottom() {
        this.messagesDiv.scrollTop = this.messagesDiv.scrollHeight;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Start the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.chatApp = new SimpleChatApp();
});

// Example function to send post comment (can be called from other pages)
function sendPostComment(postId, comment) {
    if (window.chatApp && window.chatApp.isConnected) {
        window.chatApp.sendEvent({
            type: 'POST_COMMENT',
            post_id: postId,
            comment: comment,
            user: window.chatApp.username
        });
        return true;
    }
    return false;
}
