// Simple Chat App - Fixed Version
class SimpleChatApp {
    constructor() {
        console.log('🚀 SimpleChatApp initializing...');
        
        this.ws = null;
        this.username = '';
        this.currentRoom = 'general';
        this.isConnected = false;
        
        // Wait for DOM to be ready
        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => this.init());
        } else {
            this.init();
        }
    }
    
    init() {
        console.log('📋 Initializing chat app...');
        
        this.initializeElements();
        this.bindEvents();
        this.showUsernameModal();
    }
    
    initializeElements() {
        // Get all required elements
        this.elements = {
            usernameModal: document.getElementById('usernameModal'),
            usernameInput: document.getElementById('usernameInput'),
            usernameSubmit: document.getElementById('usernameSubmit'),
            roomSelect: document.getElementById('roomSelect'),
            messagesDiv: document.getElementById('messages'),
            messageInput: document.getElementById('messageInput'),
            sendButton: document.getElementById('sendButton'),
            connectionStatus: document.getElementById('connectionStatus'),
            messageForm: document.getElementById('messageForm'),
            usernameSpan: document.getElementById('username'),
            roomNameSpan: document.getElementById('roomName')
        };
        
        // Check for missing critical elements
        const critical = ['usernameModal', 'usernameInput', 'usernameSubmit', 'messagesDiv', 'messageInput', 'connectionStatus'];
        for (const elementName of critical) {
            if (!this.elements[elementName]) {
                console.error(`❌ Critical element missing: ${elementName}`);
            } else {
                console.log(`✅ Found element: ${elementName}`);
            }
        }
    }
    
    bindEvents() {
        console.log('🔗 Binding events...');
        
        // Username modal events
        if (this.elements.usernameSubmit) {
            this.elements.usernameSubmit.addEventListener('click', () => this.setUsername());
        }
        
        if (this.elements.usernameInput) {
            this.elements.usernameInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') this.setUsername();
            });
        }
        
        // Room selection
        if (this.elements.roomSelect) {
            this.elements.roomSelect.addEventListener('change', (e) => {
                this.joinRoom(e.target.value);
            });
        }
        
        // Message form
        if (this.elements.messageForm) {
            this.elements.messageForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.sendMessage();
            });
        }
        
        if (this.elements.sendButton) {
            this.elements.sendButton.addEventListener('click', (e) => {
                e.preventDefault();
                this.sendMessage();
            });
        }
        
        if (this.elements.messageInput) {
            this.elements.messageInput.addEventListener('keypress', (e) => {
                if (e.key === 'Enter') {
                    e.preventDefault();
                    this.sendMessage();
                }
            });
        }
    }
    
    showUsernameModal() {
        console.log('👤 Showing username modal...');
        if (this.elements.usernameModal) {
            this.elements.usernameModal.style.display = 'flex';
        }
        if (this.elements.usernameInput) {
            this.elements.usernameInput.focus();
        }
    }
    
    setUsername() {
        const input = this.elements.usernameInput;
        if (!input) return;
        
        const username = input.value.trim();
        if (!username) {
            alert('Please enter a username');
            return;
        }
        
        console.log(`👤 Setting username: ${username}`);
        this.username = username;
        
        // Hide modal
        if (this.elements.usernameModal) {
            this.elements.usernameModal.style.display = 'none';
        }
        
        // Update UI
        if (this.elements.usernameSpan) {
            this.elements.usernameSpan.textContent = username;
        }
        
        // Connect to WebSocket
        this.connect();
    }
    
    connect() {
        console.log('🔌 Connecting to WebSocket...');
        this.updateStatus('connecting', '🔄 Connecting...');
        
        const wsUrl = `ws://localhost:8080/ws?username=${encodeURIComponent(this.username)}`;
        console.log(`📡 WebSocket URL: ${wsUrl}`);
        
        try {
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = () => {
                console.log('✅ WebSocket connected!');
                this.isConnected = true;
                this.updateStatus('connected', '✅ Connected');
                
                // Auto-join general room
                setTimeout(() => this.joinRoom('general'), 100);
            };
            
            this.ws.onmessage = (event) => {
                console.log('📨 Received:', event.data);
                try {
                    const data = JSON.parse(event.data);
                    this.handleMessage(data);
                } catch (error) {
                    console.error('❌ Error parsing message:', error);
                }
            };
            
            this.ws.onclose = (event) => {
                console.log('❌ WebSocket closed:', event.code, event.reason);
                this.isConnected = false;
                this.updateStatus('disconnected', '❌ Disconnected');
                
                // Try to reconnect after 3 seconds
                setTimeout(() => {
                    if (!this.isConnected) {
                        console.log('🔄 Attempting to reconnect...');
                        this.connect();
                    }
                }, 3000);
            };
            
            this.ws.onerror = (error) => {
                console.error('🚨 WebSocket error:', error);
                this.updateStatus('error', '⚠️ Connection Error');
            };
            
        } catch (error) {
            console.error('❌ Failed to create WebSocket:', error);
            this.updateStatus('error', '❌ Failed to connect');
        }
    }
    
    joinRoom(roomName) {
        if (!this.isConnected || !roomName) {
            console.log('❌ Cannot join room - not connected or no room name');
            return;
        }
        
        console.log(`🏠 Joining room: ${roomName}`);
        this.currentRoom = roomName;
        
        // Update UI
        if (this.elements.roomSelect) {
            this.elements.roomSelect.value = roomName;
        }
        if (this.elements.roomNameSpan) {
            this.elements.roomNameSpan.textContent = roomName;
        }
        
        // Clear messages
        if (this.elements.messagesDiv) {
            this.elements.messagesDiv.innerHTML = '';
        }
        
        // Send join room event
        this.sendEvent({
            type: 'JOIN_ROOM',
            room: roomName,
            user: this.username
        });
    }
    
    sendMessage() {
        const input = this.elements.messageInput;
        if (!input) return;
        
        const message = input.value.trim();
        if (!message || !this.isConnected || !this.currentRoom) {
            console.log('❌ Cannot send message:', { message: !!message, connected: this.isConnected, room: this.currentRoom });
            return;
        }
        
        console.log(`💬 Sending message: ${message}`);
        
        this.sendEvent({
            type: 'CHAT_MESSAGE',
            room: this.currentRoom,
            message: message,
            user: this.username
        });
        
        input.value = '';
    }
    
    sendEvent(event) {
        if (!this.isConnected || !this.ws) {
            console.warn('⚠️ Cannot send event - not connected');
            return false;
        }
        
        try {
            console.log('📤 Sending event:', event);
            this.ws.send(JSON.stringify(event));
            return true;
        } catch (error) {
            console.error('❌ Error sending event:', error);
            return false;
        }
    }
    
    handleMessage(data) {
        console.log('📋 Handling message:', data);
        
        switch (data.type) {
            case 'ROOM_JOINED':
                this.addSystemMessage(`✅ Joined room: ${data.room}`);
                break;
                
            case 'CHAT_MESSAGE':
                if (data.room === this.currentRoom) {
                    this.addChatMessage(data.user, data.message);
                }
                break;
                
            case 'ERROR':
                console.error('❌ Server error:', data.message);
                this.addSystemMessage(`❌ Error: ${data.message}`);
                break;
                
            default:
                console.warn('⚠️ Unknown message type:', data.type);
        }
    }
    
    addChatMessage(username, message) {
        if (!this.elements.messagesDiv) return;
        
        const messageElement = document.createElement('div');
        messageElement.className = 'message';
        
        if (username === this.username) {
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
        
        this.elements.messagesDiv.appendChild(messageElement);
        this.scrollToBottom();
    }
    
    addSystemMessage(message) {
        if (!this.elements.messagesDiv) return;
        
        const messageElement = document.createElement('div');
        messageElement.className = 'message system-message';
        messageElement.textContent = message;
        
        this.elements.messagesDiv.appendChild(messageElement);
        this.scrollToBottom();
    }
    
    updateStatus(status, text) {
        if (!this.elements.connectionStatus) return;
        
        this.elements.connectionStatus.className = `connection-status ${status}`;
        this.elements.connectionStatus.textContent = text;
    }
    
    scrollToBottom() {
        if (this.elements.messagesDiv) {
            this.elements.messagesDiv.scrollTop = this.elements.messagesDiv.scrollHeight;
        }
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Initialize when DOM is ready
console.log('📦 Loading SimpleChatApp...');
window.chatApp = new SimpleChatApp();
