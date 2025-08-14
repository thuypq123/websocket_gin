// WebSocket Adapter - Clean separation between WebSocket and UI logic
class WebSocketAdapter {
    constructor(config) {
        this.config = {
            url: config.url,
            username: config.username || 'Anonymous',
            userId: config.userId || this.generateId(),
            room: config.room || 'general',
            roomType: config.roomType || 'chat',
            reconnectInterval: config.reconnectInterval || 3000,
            maxReconnectAttempts: config.maxReconnectAttempts || 5,
            ...config
        };
        
        this.ws = null;
        this.isConnected = false;
        this.reconnectAttempts = 0;
        this.eventHandlers = new Map();
        this.connectionStateHandlers = [];
        
        this.setupDefaultHandlers();
    }
    
    // Connection Management
    connect() {
        const wsUrl = this.buildWebSocketUrl();
        console.log(`🔌 Connecting to WebSocket: ${wsUrl}`);
        
        this.ws = new WebSocket(wsUrl);
        this.setupWebSocketEventHandlers();
        
        this.notifyConnectionState('connecting');
    }
    
    disconnect() {
        if (this.ws) {
            this.isConnected = false;
            this.ws.close();
        }
    }
    
    // Event System
    on(eventType, handler) {
        if (!this.eventHandlers.has(eventType)) {
            this.eventHandlers.set(eventType, []);
        }
        this.eventHandlers.get(eventType).push(handler);
    }
    
    off(eventType, handler) {
        if (this.eventHandlers.has(eventType)) {
            const handlers = this.eventHandlers.get(eventType);
            const index = handlers.indexOf(handler);
            if (index > -1) {
                handlers.splice(index, 1);
            }
        }
    }
    
    emit(eventType, data) {
        if (this.eventHandlers.has(eventType)) {
            this.eventHandlers.get(eventType).forEach(handler => {
                try {
                    handler(data);
                } catch (error) {
                    console.error(`Error in event handler for ${eventType}:`, error);
                }
            });
        }
    }
    
    // Connection State Management
    onConnectionStateChange(handler) {
        this.connectionStateHandlers.push(handler);
    }
    
    notifyConnectionState(state) {
        this.connectionStateHandlers.forEach(handler => {
            try {
                handler(state, this.isConnected);
            } catch (error) {
                console.error('Error in connection state handler:', error);
            }
        });
    }
    
    // Message Sending
    sendEvent(eventType, action, data) {
        if (!this.isConnected) {
            console.warn('⚠️ Cannot send event: WebSocket not connected');
            return false;
        }
        
        const wsEvent = {
            type: eventType,
            action: action,
            data: data,
            timestamp: new Date().toISOString()
        };
        
        try {
            this.ws.send(JSON.stringify(wsEvent));
            console.log(`📤 Sent event: ${eventType}.${action}`, wsEvent);
            return true;
        } catch (error) {
            console.error('❌ Failed to send event:', error);
            return false;
        }
    }
    
    // Convenience methods for specific event types
    sendChatMessage(content) {
        return this.sendEvent('chat', 'send', {
            message: { content: content }
        });
    }
    
    sendComment(postId, content) {
        return this.sendEvent('comment', 'create', {
            comment: { content: content },
            post_id: postId
        });
    }
    
    updateComment(commentId, content) {
        return this.sendEvent('comment', 'update', {
            comment: { id: commentId, content: content }
        });
    }
    
    deleteComment(commentId) {
        return this.sendEvent('comment', 'delete', {
            comment: { id: commentId }
        });
    }
    
    // ===== SUBSCRIPTION METHODS =====
    
    /**
     * Subscribe to specific events
     * @param {string} eventType - Event type (chat, comment, post)
     * @param {string} eventAction - Event action (send, create, update, delete, *)
     * @param {string} resourceId - Resource ID or * for all
     * @param {Object} filters - Additional filters
     * @returns {Promise<string>} - Subscription ID
     */
    async subscribe(eventType, eventAction, resourceId, filters = {}) {
        return new Promise((resolve, reject) => {
            const subscriptionRequest = {
                action: 'subscribe',
                event_type: eventType,
                event_action: eventAction,
                resource_id: resourceId,
                filters: filters
            };
            
            // Set up one-time response handler
            const responseHandler = (data) => {
                if (data.success !== undefined) {
                    this.off('message', responseHandler);
                    if (data.success) {
                        console.log(`✅ Subscribed: ${eventType}.${eventAction}.${resourceId}`);
                        resolve(data.subscription_id);
                    } else {
                        reject(new Error(data.error || 'Subscription failed'));
                    }
                }
            };
            
            this.on('message', responseHandler);
            
            // Send subscription request as raw JSON
            if (!this.sendRawMessage(JSON.stringify(subscriptionRequest))) {
                this.off('message', responseHandler);
                reject(new Error('Failed to send subscription request'));
            }
            
            // Timeout after 5 seconds
            setTimeout(() => {
                this.off('message', responseHandler);
                reject(new Error('Subscription request timeout'));
            }, 5000);
        });
    }
    
    /**
     * Subscribe to comments on a specific post
     */
    async subscribeToPostComments(postId) {
        return this.subscribe('comment', '*', postId);
    }
    
    /**
     * Subscribe to messages in a specific chat room
     */
    async subscribeToChatRoom(roomId) {
        return this.subscribe('chat', 'send', roomId);
    }
    
    /**
     * Unsubscribe from all subscriptions
     */
    async unsubscribeAll() {
        return new Promise((resolve, reject) => {
            const unsubscribeRequest = {
                action: 'unsubscribe'
            };
            
            const responseHandler = (data) => {
                if (data.success !== undefined) {
                    this.off('message', responseHandler);
                    if (data.success) {
                        console.log('❌ Unsubscribed from all subscriptions');
                        resolve(true);
                    } else {
                        reject(new Error(data.error || 'Unsubscribe failed'));
                    }
                }
            };
            
            this.on('message', responseHandler);
            
            if (!this.sendRawMessage(JSON.stringify(unsubscribeRequest))) {
                this.off('message', responseHandler);
                reject(new Error('Failed to send unsubscribe request'));
            }
            
            setTimeout(() => {
                this.off('message', responseHandler);
                reject(new Error('Unsubscribe request timeout'));
            }, 5000);
        });
    }
    
    /**
     * Send raw message (for subscription requests)
     */
    sendRawMessage(message) {
        if (!this.isConnected) {
            console.warn('⚠️ Cannot send message: WebSocket not connected');
            return false;
        }
        
        try {
            this.ws.send(message);
            return true;
        } catch (error) {
            console.error('❌ Failed to send raw message:', error);
            return false;
        }
    }
    
    // Private Methods
    buildWebSocketUrl() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const params = new URLSearchParams({
            username: this.config.username,
            user_id: this.config.userId,
            room: this.config.room,
            room_type: this.config.roomType
        });
        
        return `${protocol}//${window.location.host}/ws?${params.toString()}`;
    }
    
    setupWebSocketEventHandlers() {
        this.ws.onopen = () => {
            console.log('✅ WebSocket connected');
            this.isConnected = true;
            this.reconnectAttempts = 0;
            this.notifyConnectionState('connected');
            this.emit('connected', {});
        };
        
        this.ws.onmessage = (event) => {
            try {
                const wsEvent = JSON.parse(event.data);
                console.log(`📥 Received event: ${wsEvent.type}.${wsEvent.action}`, wsEvent);
                
                // Emit specific event type
                this.emit(wsEvent.type, wsEvent);
                
                // Emit general message event
                this.emit('message', wsEvent);
                
                // Emit specific action events
                this.emit(`${wsEvent.type}.${wsEvent.action}`, wsEvent);
                
            } catch (error) {
                console.error('❌ Error parsing WebSocket message:', error);
                this.emit('error', { type: 'parse_error', error: error });
            }
        };
        
        this.ws.onclose = (event) => {
            console.log('🔌 WebSocket connection closed', event);
            this.isConnected = false;
            this.notifyConnectionState('disconnected');
            this.emit('disconnected', { event: event });
            
            // Auto-reconnect logic
            if (this.reconnectAttempts < this.config.maxReconnectAttempts) {
                this.scheduleReconnect();
            } else {
                console.error('❌ Max reconnect attempts reached');
                this.emit('reconnect_failed', {});
            }
        };
        
        this.ws.onerror = (error) => {
            console.error('❌ WebSocket error:', error);
            this.notifyConnectionState('error');
            this.emit('error', { type: 'connection_error', error: error });
        };
    }
    
    scheduleReconnect() {
        this.reconnectAttempts++;
        console.log(`🔄 Scheduling reconnect attempt ${this.reconnectAttempts}/${this.config.maxReconnectAttempts}`);
        
        this.notifyConnectionState('reconnecting');
        
        setTimeout(() => {
            if (!this.isConnected) {
                this.connect();
            }
        }, this.config.reconnectInterval);
    }
    
    setupDefaultHandlers() {
        // Default error handler
        this.on('error', (data) => {
            console.error('WebSocket Adapter Error:', data);
        });
        
        // Default connection handlers
        this.on('connected', () => {
            console.log('🎉 WebSocket Adapter connected successfully');
        });
        
        this.on('disconnected', () => {
            console.log('👋 WebSocket Adapter disconnected');
        });
    }
    
    generateId() {
        return 'user_' + Math.random().toString(36).substr(2, 9);
    }
    
    // Utility Methods
    getConnectionState() {
        if (!this.ws) return 'not_initialized';
        
        switch (this.ws.readyState) {
            case WebSocket.CONNECTING: return 'connecting';
            case WebSocket.OPEN: return 'connected';
            case WebSocket.CLOSING: return 'closing';
            case WebSocket.CLOSED: return 'disconnected';
            default: return 'unknown';
        }
    }
    
    getConfig() {
        return { ...this.config };
    }
    
    updateConfig(newConfig) {
        this.config = { ...this.config, ...newConfig };
    }
}

// Export for different module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = WebSocketAdapter;
} else if (typeof window !== 'undefined') {
    window.WebSocketAdapter = WebSocketAdapter;
}
