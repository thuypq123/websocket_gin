// Multi-Channel WebSocket Demo
class MultiChannelDemo {
    constructor() {
        this.wsAdapter = null;
        this.username = '';
        this.subscriptions = new Set();
        
        this.init();
    }
    
    init() {
        this.setupUI();
        this.bindEvents();
        this.showUsernameModal();
    }
    
    setupUI() {
        // Get DOM elements
        this.elements = {
            usernameModal: document.getElementById('usernameModal'),
            usernameInput: document.getElementById('usernameInput'),
            usernameSubmit: document.getElementById('usernameSubmit'),
            
            connectionStatus: document.getElementById('connectionStatus'),
            currentUser: document.getElementById('currentUser'),
            
            subscriptionsList: document.getElementById('subscriptionsList'),
            channelInput: document.getElementById('channelInput'),
            subscribeBtn: document.getElementById('subscribeBtn'),
            unsubscribeBtn: document.getElementById('unsubscribeBtn'),
            listSubscriptionsBtn: document.getElementById('listSubscriptionsBtn'),
            
            messageArea: document.getElementById('messageArea'),
            messageInput: document.getElementById('messageInput'),
            sendMessageBtn: document.getElementById('sendMessageBtn'),
            
            quickChannels: document.getElementById('quickChannels'),
            statsBtn: document.getElementById('statsBtn'),
            statsArea: document.getElementById('statsArea')
        };
    }
    
    bindEvents() {
        // Username modal
        this.elements.usernameSubmit?.addEventListener('click', () => this.setUsername());
        this.elements.usernameInput?.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.setUsername();
        });
        
        // Subscription management
        this.elements.subscribeBtn?.addEventListener('click', () => this.subscribeToChannel());
        this.elements.unsubscribeBtn?.addEventListener('click', () => this.unsubscribeFromChannel());
        this.elements.listSubscriptionsBtn?.addEventListener('click', () => this.listSubscriptions());
        
        this.elements.channelInput?.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.subscribeToChannel();
        });
        
        // Message sending
        this.elements.sendMessageBtn?.addEventListener('click', () => this.sendMessage());
        this.elements.messageInput?.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.sendMessage();
        });
        
        // Quick channel buttons
        this.elements.quickChannels?.addEventListener('click', (e) => {
            if (e.target.classList.contains('quick-channel-btn')) {
                const channel = e.target.dataset.channel;
                this.subscribeToChannels([channel]);
            }
        });
        
        // Stats
        this.elements.statsBtn?.addEventListener('click', () => this.loadStats());
    }
    
    showUsernameModal() {
        if (this.elements.usernameModal) {
            this.elements.usernameModal.style.display = 'flex';
            this.elements.usernameInput?.focus();
        }
    }
    
    setUsername() {
        const username = this.elements.usernameInput?.value.trim();
        if (!username) {
            alert('Please enter a username');
            return;
        }
        
        this.username = username;
        this.elements.usernameModal.style.display = 'none';
        this.elements.currentUser.textContent = username;
        
        this.setupWebSocket();
    }
    
    setupWebSocket() {
        this.wsAdapter = new WebSocketAdapter({
            url: `ws://localhost:8080/ws`,
            username: this.username
        });
        
        // Connection state handlers
        this.wsAdapter.onConnectionStateChange((state, isConnected) => {
            this.updateConnectionStatus(state, isConnected);
        });
        
        // Event handlers
        this.wsAdapter.on('SUBSCRIPTION_RESPONSE', (data) => this.handleSubscriptionResponse(data));
        this.wsAdapter.on('CHAT_MESSAGE', (data) => this.handleChatMessage(data));
        this.wsAdapter.on('POST_COMMENT', (data) => this.handlePostComment(data));
        this.wsAdapter.on('ROOM_JOINED', (data) => this.handleRoomJoined(data));
        this.wsAdapter.on('ERROR', (data) => this.handleError(data));
        
        this.wsAdapter.connect();
    }
    
    updateConnectionStatus(state, isConnected) {
        const status = this.elements.connectionStatus;
        if (!status) return;
        
        status.textContent = isConnected ? 'Connected' : `Connecting... (${state})`;
        status.className = `status ${isConnected ? 'connected' : 'connecting'}`;
    }
    
    subscribeToChannel() {
        const channel = this.elements.channelInput?.value.trim();
        if (!channel) {
            alert('Please enter a channel name');
            return;
        }
        
        this.subscribeToChannels([channel]);
        this.elements.channelInput.value = '';
    }
    
    subscribeToChannels(channels) {
        if (!this.wsAdapter) return;
        
        this.wsAdapter.subscribe(channels);
        this.addMessage('System', `Subscribing to: ${channels.join(', ')}`, 'system');
    }
    
    unsubscribeFromChannel() {
        const channel = this.elements.channelInput?.value.trim();
        if (!channel) {
            alert('Please enter a channel name');
            return;
        }
        
        this.wsAdapter.unsubscribe([channel]);
        this.addMessage('System', `Unsubscribing from: ${channel}`, 'system');
        this.elements.channelInput.value = '';
    }
    
    listSubscriptions() {
        if (!this.wsAdapter) return;
        
        this.wsAdapter.listSubscriptions();
        this.addMessage('System', 'Requesting subscription list...', 'system');
    }
    
    sendMessage() {
        const message = this.elements.messageInput?.value.trim();
        if (!message) return;
        
        // Parse message format: /channel_name message or just message (defaults to chat:general)
        let channelName = 'chat:general';
        let messageContent = message;
        
        if (message.startsWith('/')) {
            const parts = message.split(' ');
            if (parts.length > 1) {
                channelName = parts[0].substring(1); // Remove /
                messageContent = parts.slice(1).join(' ');
            }
        }
        
        this.sendToChannel(channelName, messageContent);
        this.elements.messageInput.value = '';
    }
    
    sendToChannel(channelName, message) {
        if (!this.wsAdapter) return;
        
        // Parse channel to determine message type
        const [channelType, identifier] = channelName.split(':', 2);
        
        switch (channelType) {
            case 'chat':
                this.wsAdapter.sendEvent('CHAT_MESSAGE', 'send', {
                    type: 'CHAT_MESSAGE',
                    room: identifier,
                    message: message,
                    user: this.username
                });
                break;
            case 'post':
                this.wsAdapter.sendEvent('POST_COMMENT', 'send', {
                    type: 'POST_COMMENT',
                    post_id: identifier,
                    comment: message,
                    user: this.username
                });
                break;
            default:
                this.addMessage('System', `Unknown channel type: ${channelType}`, 'error');
                return;
        }
        
        this.addMessage('You', `[${channelName}] ${message}`, 'sent');
    }
    
    handleSubscriptionResponse(data) {
        switch (data.action) {
            case 'subscribed':
                data.channels.forEach(channel => this.subscriptions.add(channel));
                this.addMessage('System', `✅ Subscribed to: ${data.channels.join(', ')}`, 'success');
                break;
            case 'unsubscribed':
                data.channels.forEach(channel => this.subscriptions.delete(channel));
                this.addMessage('System', `❌ Unsubscribed from: ${data.channels.join(', ')}`, 'success');
                break;
            case 'listed':
                this.subscriptions = new Set(data.subscriptions);
                this.addMessage('System', `📋 Current subscriptions: ${data.subscriptions.join(', ') || 'None'}`, 'info');
                break;
        }
        
        this.updateSubscriptionsList();
        
        if (!data.success && data.message) {
            this.addMessage('System', `⚠️ ${data.message}`, 'warning');
        }
    }
    
    handleChatMessage(data) {
        this.addMessage(data.user, `[chat:${data.room}] ${data.message}`, 'received');
    }
    
    handlePostComment(data) {
        this.addMessage(data.user, `[post:${data.post_id}] ${data.comment}`, 'received');
    }
    
    handleRoomJoined(data) {
        this.addMessage('System', `🏠 Joined room: ${data.room}`, 'info');
    }
    
    handleError(data) {
        this.addMessage('Error', data.message, 'error');
    }
    
    addMessage(sender, message, type = 'default') {
        if (!this.elements.messageArea) return;
        
        const messageDiv = document.createElement('div');
        messageDiv.className = `message ${type}`;
        messageDiv.innerHTML = `
            <span class="sender">${sender}:</span>
            <span class="content">${message}</span>
            <span class="time">${new Date().toLocaleTimeString()}</span>
        `;
        
        this.elements.messageArea.appendChild(messageDiv);
        this.elements.messageArea.scrollTop = this.elements.messageArea.scrollHeight;
    }
    
    updateSubscriptionsList() {
        if (!this.elements.subscriptionsList) return;
        
        const subscriptions = Array.from(this.subscriptions);
        this.elements.subscriptionsList.innerHTML = subscriptions.length > 0 
            ? subscriptions.map(sub => `<span class="subscription-tag">${sub}</span>`).join(' ')
            : '<span class="no-subscriptions">No active subscriptions</span>';
    }
    
    async loadStats() {
        try {
            const response = await fetch('/api/v1/ws/stats');
            const stats = await response.json();
            
            this.elements.statsArea.innerHTML = `
                <h3>WebSocket Statistics</h3>
                <div class="stats-grid">
                    <div class="stat-item">
                        <strong>Total Clients:</strong> ${stats.websocket_stats.total_clients}
                    </div>
                    <div class="stat-item">
                        <strong>Active Channels:</strong> ${Object.keys(stats.channel_stats || {}).length}
                    </div>
                    <div class="stat-item">
                        <strong>Channel Details:</strong>
                        <pre>${JSON.stringify(stats.channel_stats, null, 2)}</pre>
                    </div>
                </div>
            `;
        } catch (error) {
            this.addMessage('System', `Failed to load stats: ${error.message}`, 'error');
        }
    }
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new MultiChannelDemo();
});
