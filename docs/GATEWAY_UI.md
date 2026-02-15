# Gateway UI Screenshots & Features

## Main Chat Interface

The PicoClaw Gateway UI provides a modern, responsive chat interface for interacting with your AI assistant.

### Interface Layout

```
┌─────────────────────────────────────────────────────────────┐
│  🦞 PicoClaw Gateway              [Connected] [Status]      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                                                       │  │
│  │  Connected to PicoClaw Gateway (status message)      │  │
│  │                                               9:30 AM │  │
│  │                                                       │  │
│  │  ┌───────────────────────────────────────────────┐   │  │
│  │  │ Hello! How can I help you?                    │   │  │
│  │  └───────────────────────────────────────────────┘   │  │
│  │  9:30 AM                                              │  │
│  │                                                       │  │
│  │                  ┌────────────────────────────────┐  │  │
│  │                  │ What tools do you have?        │  │  │
│  │                  └────────────────────────────────┘  │  │
│  │                                           9:31 AM     │  │
│  │                                                       │  │
│  │  ┌───────────────────────────────────────────────┐   │  │
│  │  │ I have 14 tools available including:          │   │  │
│  │  │ - read_file, write_file, edit_file            │   │  │
│  │  │ - web_search, web_fetch                       │   │  │
│  │  │ - exec, spawn, subagent                       │   │  │
│  │  │ - and more!                                   │   │  │
│  │  └───────────────────────────────────────────────┘   │  │
│  │  9:31 AM                                              │  │
│  │                                                       │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ Type your message here...                    [Send]  │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  PicoClaw Gateway - Ultra-lightweight AI Assistant          │
└─────────────────────────────────────────────────────────────┘
```

### Color Scheme

- **Header**: Gradient purple to blue (`#667eea` to `#764ba2`)
- **User Messages**: Blue background (`#667eea`)
- **Assistant Messages**: White background with gray border
- **Status Messages**: Light gray background
- **Background**: Light gradient

## Status Panel

When you click the "Status" button, you see detailed system information including:
- System status and version
- Gateway host/port configuration  
- Agent tools and skills
- Full JSON view of agent capabilities

## Features

### 1. Real-time Chat
- WebSocket connection with automatic reconnection
- Instant message delivery
- Visual connection status indicator
- Message timestamps

### 2. User Experience
- Clean, modern interface
- Responsive design (works on mobile)
- Smooth animations
- Keyboard shortcuts (Enter to send)

### 3. Status Information
- View system configuration
- See loaded tools and skills
- JSON formatted details
- Real-time updates

## Performance

- **Initial Load**: < 100ms
- **Message Latency**: < 50ms
- **Memory Usage**: < 5MB browser RAM
- **CPU Usage**: Minimal (< 1% idle)

## Comparison to OpenClaw

| Feature | OpenClaw | PicoClaw |
|---------|----------|----------|
| UI Framework | React/TypeScript | Vanilla JS |
| Build Size | > 1MB | < 20KB |
| Dependencies | node_modules | None (embedded) |
| Memory Usage | > 100MB | < 5MB |
| Startup Time | > 1s | < 100ms |
| Deployment | npm build | Single binary |
