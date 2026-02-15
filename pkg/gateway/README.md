# PicoClaw Gateway UI

The PicoClaw Gateway UI provides a web-based interface for interacting with your PicoClaw AI assistant, similar to OpenClaw's interface.

## Features

- **Real-time Chat Interface**: WebSocket-based communication for instant messaging with your AI assistant
- **Status Dashboard**: View agent configuration, loaded tools, and available skills
- **Responsive Design**: Works on desktop and mobile devices
- **Ultra-Lightweight**: Minimal resource usage, keeping with PicoClaw's philosophy

## Architecture

The gateway UI consists of:

1. **HTTP Server** (`server.go`): 
   - Serves static files (HTML, CSS, JS)
   - Provides WebSocket endpoint for chat
   - Exposes REST API for status information

2. **Web Interface** (`ui/`):
   - `index.html`: Main chat interface
   - `style.css`: Modern, responsive styling
   - `app.js`: WebSocket client and UI logic

## Usage

The gateway UI is automatically started when you run the gateway command:

```bash
picoclaw gateway
```

Then open your browser to:

```
http://localhost:18790
```

Or configure a different host/port in your config:

```json
{
  "gateway": {
    "host": "0.0.0.0",
    "port": 18790
  }
}
```

## API Endpoints

### WebSocket Chat

- **Endpoint**: `ws://localhost:18790/ws/chat`
- **Protocol**: JSON messages
- **Message Format**:
  ```json
  {
    "type": "user",
    "message": "Hello, PicoClaw!",
    "timestamp": 1234567890
  }
  ```

### Status API

- **Endpoint**: `GET /api/status`
- **Response**:
  ```json
  {
    "version": "1.0.0",
    "status": "running",
    "agent": {
      "tools": {
        "count": 14,
        "names": ["read_file", "write_file", ...]
      },
      "skills": {
        "total": 0,
        "available": 0,
        "names": []
      }
    },
    "gateway": {
      "host": "0.0.0.0",
      "port": 18790
    }
  }
  ```

## UI Features

### Chat Interface

The chat interface allows you to:
- Send messages to your AI assistant
- View conversation history
- See connection status
- Real-time message delivery

### Status Panel

Click the "Status" button to view:
- System status and version
- Gateway configuration
- Agent information (tools, skills)
- Detailed JSON view of the system state

## Security Considerations

The gateway UI currently:
- Accepts connections from all origins (for simplicity)
- Does not implement authentication
- Should only be exposed on trusted networks

For production use, consider:
- Implementing authentication/authorization
- Using HTTPS with TLS
- Restricting allowed origins
- Adding rate limiting

## Development

To modify the UI:

1. Edit files in `pkg/gateway/ui/`
2. Files are embedded into the binary at build time
3. Rebuild with `make build`

The UI files are embedded using Go's `embed` directive, so they're part of the single binary.

## Comparison to OpenClaw

Like OpenClaw's gateway, this provides:
- Web-based chat interface
- Real-time communication
- Status and monitoring

Unlike OpenClaw, PicoClaw's gateway:
- Uses <10MB RAM (vs >1GB for OpenClaw)
- Single binary deployment
- Embedded static files
- Go-based for better performance
