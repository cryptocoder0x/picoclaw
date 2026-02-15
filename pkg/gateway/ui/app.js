// WebSocket connection
let ws = null;
let reconnectTimer = null;

// DOM elements
const messagesContainer = document.getElementById('messages');
const messageInput = document.getElementById('message-input');
const sendBtn = document.getElementById('send-btn');
const connectionStatus = document.getElementById('connection-status');
const statusBtn = document.getElementById('status-btn');
const statusPanel = document.getElementById('status-panel');
const statusContent = document.getElementById('status-content');
const closeStatusBtn = document.getElementById('close-status-btn');

// Initialize
function init() {
    connectWebSocket();
    setupEventListeners();
}

// Connect to WebSocket
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws/chat`;

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('WebSocket connected');
        updateConnectionStatus(true);
    };

    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        displayMessage(message);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };

    ws.onclose = () => {
        console.log('WebSocket disconnected');
        updateConnectionStatus(false);
        // Attempt to reconnect after 3 seconds
        reconnectTimer = setTimeout(() => {
            connectWebSocket();
        }, 3000);
    };
}

// Update connection status indicator
function updateConnectionStatus(connected) {
    if (connected) {
        connectionStatus.textContent = 'Connected';
        connectionStatus.className = 'status-indicator connected';
        sendBtn.disabled = false;
        messageInput.disabled = false;
    } else {
        connectionStatus.textContent = 'Disconnected';
        connectionStatus.className = 'status-indicator disconnected';
        sendBtn.disabled = true;
        messageInput.disabled = true;
    }
}

// Display a message in the chat
function displayMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${message.type}`;

    const contentDiv = document.createElement('div');
    contentDiv.className = 'message-content';
    contentDiv.textContent = message.message;

    const timeDiv = document.createElement('div');
    timeDiv.className = 'message-time';
    timeDiv.textContent = formatTime(message.timestamp);

    messageDiv.appendChild(contentDiv);
    messageDiv.appendChild(timeDiv);
    messagesContainer.appendChild(messageDiv);

    // Scroll to bottom
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

// Send a message
function sendMessage() {
    const text = messageInput.value.trim();
    if (!text || !ws || ws.readyState !== WebSocket.OPEN) {
        return;
    }

    const message = {
        type: 'user',
        message: text,
        timestamp: Math.floor(Date.now() / 1000)
    };

    // Display user message immediately
    displayMessage(message);

    // Send to server
    ws.send(JSON.stringify(message));

    // Clear input
    messageInput.value = '';
}

// Fetch and display status
async function fetchStatus() {
    try {
        const response = await fetch('/api/status');
        const data = await response.json();
        
        statusContent.innerHTML = `
            <div style="margin-bottom: 20px;">
                <h3 style="margin-bottom: 10px;">System Status</h3>
                <p><strong>Status:</strong> ${data.status}</p>
                <p><strong>Version:</strong> ${data.version}</p>
            </div>
            <div style="margin-bottom: 20px;">
                <h3 style="margin-bottom: 10px;">Gateway Configuration</h3>
                <p><strong>Host:</strong> ${data.gateway.host}</p>
                <p><strong>Port:</strong> ${data.gateway.port}</p>
            </div>
            <div>
                <h3 style="margin-bottom: 10px;">Agent Information</h3>
                <pre>${JSON.stringify(data.agent, null, 2)}</pre>
            </div>
        `;
    } catch (error) {
        statusContent.innerHTML = `<p style="color: #ef4444;">Error fetching status: ${error.message}</p>`;
    }
}

// Setup event listeners
function setupEventListeners() {
    sendBtn.addEventListener('click', sendMessage);
    
    messageInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });

    statusBtn.addEventListener('click', () => {
        statusPanel.classList.remove('hidden');
        fetchStatus();
    });

    closeStatusBtn.addEventListener('click', () => {
        statusPanel.classList.add('hidden');
    });
}

// Format timestamp
function formatTime(timestamp) {
    const date = new Date(timestamp * 1000);
    return date.toLocaleTimeString();
}

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    if (reconnectTimer) {
        clearTimeout(reconnectTimer);
    }
    if (ws) {
        ws.close();
    }
});

// Start the application
init();
