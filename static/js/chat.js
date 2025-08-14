class ChatApp {
    constructor() {
        this.wsAdapter = null;
        this.username = '';
        this.room = '';
        this.isConnected = false;
        this.subscriptionId = null;
        
        this.initializeElements();
        this.parseParams();
        this.setupEventListeners();
        this.loadMessageHistory();
        this.connect();
    }

    initializeElements() {
        this.messagesContainer = document.getElementById('messages');
        this.messageForm = document.getElementById('messageForm');
        this.messageInput = document.getElementById('messageInput');
        this.connectionStatus = document.getElementById('connectionStatus');
        this.usernameDisplay = document.getElementById('username');
        this.roomNameDisplay = document.getElementById('roomName');
        this.leaveBtn = document.getElementById('leaveBtn');
    }

    parseParams() {
        const params = new URLSearchParams(window.location.search);
        this.username = params.get('username') || 'Anonymous';
        this.room = params.get('room') || 'general';
        
        this.usernameDisplay.textContent = this.username;
        this.roomNameDisplay.textContent = this.room;
    }

    setupEventListeners() {
        this.messageForm.addEventListener('submit', (e) => {
            e.preventDefault();
            this.sendMessage();
        });

        this.leaveBtn.addEventListener('click', () => {
            this.disconnect();
            window.location.href = '/';
        });

        window.addEventListener('beforeunload', () => {
            this.disconnect();
        });
    }

    connect() {
        this.updateConnectionStatus('connecting');
        
        // Create WebSocket adapter
        this.wsAdapter = new WebSocketAdapter({
            username: this.username,
            room: this.room,
            roomType: 'chat'
        });
        
        // Set up event handlers
        this.wsAdapter.onConnectionStateChange((state, isConnected) => {
            this.isConnected = isConnected;
            this.updateConnectionStatus(state);
        });
        
        this.wsAdapter.on('connected', async () => {
            console.log('✅ Connected to WebSocket');
            
            // Subscribe to chat messages in this room
            try {
                this.subscriptionId = await this.wsAdapter.subscribeToChatRoom(this.room);
                console.log(`📋 Subscribed to chat room: ${this.room}`);
            } catch (error) {
                console.error('❌ Failed to subscribe to chat room:', error);
            }
        });
        
        // Handle all WebSocket events
        this.wsAdapter.on('message', (wsEvent) => {
            this.handleWebSocketEvent(wsEvent);
        });
        
        this.wsAdapter.on('disconnected', () => {
            console.log('❌ WebSocket disconnected');
            this.subscriptionId = null;
        });
        
        // Connect
        this.wsAdapter.connect();
    }

    disconnect() {
        if (this.wsAdapter) {
            this.isConnected = false;
            this.wsAdapter.disconnect();
        }
    }

    sendMessage() {
        const content = this.messageInput.value.trim();
        
        if (!content || !this.isConnected) {
            return;
        }

        // Use WebSocket adapter to send chat message
        if (this.wsAdapter.sendChatMessage(content)) {
            this.messageInput.value = '';
        }
    }

    handleWebSocketEvent(wsEvent) {
        // Handle different event types
        if (wsEvent.type === 'chat') {
            if (wsEvent.action === 'send' || wsEvent.action === 'join' || wsEvent.action === 'leave') {
                // Extract message data from WSEvent
                const message = {
                    id: wsEvent.event_id,
                    username: wsEvent.username,
                    content: wsEvent.data.message ? wsEvent.data.message.content : wsEvent.username + (wsEvent.action === 'join' ? ' joined the chat' : ' left the chat'),
                    room_id: wsEvent.room_id,
                    timestamp: wsEvent.timestamp,
                    type: wsEvent.action === 'send' ? 'message' : wsEvent.action
                };
                this.displayMessage(message);
            }
        } else if (wsEvent.type === 'error') {
            console.error('WebSocket error:', wsEvent.data.error);
        }
    }

    displayMessage(message) {
        const messageElement = document.createElement('div');
        messageElement.className = 'message';
        
        if (message.type === 'join' || message.type === 'leave') {
            messageElement.className += ' system';
            messageElement.innerHTML = `
                <div class="message-content">${message.content}</div>
            `;
        } else {
            const isOwnMessage = message.username === this.username;
            messageElement.className += isOwnMessage ? ' own' : ' other';
            
            const timestamp = new Date(message.timestamp).toLocaleTimeString();
            
            messageElement.innerHTML = `
                <div class="message-header">
                    ${isOwnMessage ? 'You' : message.username} • ${timestamp}
                </div>
                <div class="message-content">${this.escapeHtml(message.content)}</div>
            `;
        }

        this.messagesContainer.appendChild(messageElement);
        this.scrollToBottom();
    }

    updateConnectionStatus(status) {
        this.connectionStatus.className = `connection-status ${status}`;
        
        switch (status) {
            case 'connected':
                this.connectionStatus.textContent = 'Connected';
                break;
            case 'connecting':
                this.connectionStatus.textContent = 'Connecting...';
                break;
            case 'disconnected':
                this.connectionStatus.textContent = 'Disconnected';
                break;
        }
    }

    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    async loadMessageHistory() {
        try {
            const response = await fetch(`/api/v1/messages/recent?room=${encodeURIComponent(this.room)}&limit=20`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            const data = await response.json();
            const messages = data.messages || [];
            
            // Clear existing messages
            this.messagesContainer.innerHTML = '';
            
            // Display historical messages
            messages.forEach(message => {
                this.displayMessage(message);
            });
            
            // Add separator between history and new messages
            if (messages.length > 0) {
                const separator = document.createElement('div');
                separator.className = 'message-separator';
                separator.innerHTML = '<span>--- Recent Messages ---</span>';
                this.messagesContainer.appendChild(separator);
            }
            
            console.log(`Loaded ${messages.length} messages from history`);
            
        } catch (error) {
            console.error('Failed to load message history:', error);
            // Don't show error to user, just log it
        }
    }
}

// Initialize the chat app when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});