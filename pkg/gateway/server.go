package gateway

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sipeed/picoclaw/pkg/agent"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
)

//go:embed ui/*
var uiFiles embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

type Server struct {
	config     *config.Config
	agentLoop  *agent.AgentLoop
	httpServer *http.Server
	clients    map[*websocket.Conn]bool
	clientsMu  sync.RWMutex
}

type ChatMessage struct {
	Type      string `json:"type"`      // "user" or "assistant" or "status"
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func NewServer(cfg *config.Config, agentLoop *agent.AgentLoop) *Server {
	return &Server{
		config:    cfg,
		agentLoop: agentLoop,
		clients:   make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve static files from embedded FS
	uiFS, err := fs.Sub(uiFiles, "ui")
	if err != nil {
		return fmt.Errorf("failed to access UI files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(uiFS)))

	// WebSocket endpoint for chat
	mux.HandleFunc("/ws/chat", s.handleWebSocket)

	// API endpoints
	mux.HandleFunc("/api/status", s.handleStatus)

	addr := fmt.Sprintf("%s:%d", s.config.Gateway.Host, s.config.Gateway.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.InfoCF("gateway", "Starting web UI server", map[string]interface{}{
		"address": addr,
	})

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorCF("gateway", "HTTP server error", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ErrorCF("gateway", "WebSocket upgrade failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	defer conn.Close()

	// Register client
	s.clientsMu.Lock()
	s.clients[conn] = true
	s.clientsMu.Unlock()

	// Unregister on disconnect
	defer func() {
		s.clientsMu.Lock()
		delete(s.clients, conn)
		s.clientsMu.Unlock()
	}()

	// Send welcome message
	welcomeMsg := ChatMessage{
		Type:      "status",
		Message:   "Connected to PicoClaw Gateway",
		Timestamp: time.Now().Unix(),
	}
	if err := conn.WriteJSON(welcomeMsg); err != nil {
		logger.ErrorCF("gateway", "Failed to send welcome message", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Handle incoming messages
	for {
		var msg ChatMessage
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.ErrorCF("gateway", "WebSocket error", map[string]interface{}{
					"error": err.Error(),
				})
			}
			break
		}

		if msg.Type == "user" {
			// Process message through agent
			response, err := s.agentLoop.ProcessDirect(context.Background(), msg.Message, "web:ui")
			if err != nil {
				responseMsg := ChatMessage{
					Type:      "assistant",
					Message:   fmt.Sprintf("Error: %v", err),
					Timestamp: time.Now().Unix(),
				}
				conn.WriteJSON(responseMsg)
				continue
			}

			// Send response
			responseMsg := ChatMessage{
				Type:      "assistant",
				Message:   response,
				Timestamp: time.Now().Unix(),
			}
			if err := conn.WriteJSON(responseMsg); err != nil {
				logger.ErrorCF("gateway", "Failed to send response", map[string]interface{}{
					"error": err.Error(),
				})
				break
			}
		}
	}
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	startupInfo := s.agentLoop.GetStartupInfo()
	
	status := map[string]interface{}{
		"version":    "1.0.0",
		"status":     "running",
		"agent":      startupInfo,
		"gateway": map[string]interface{}{
			"host": s.config.Gateway.Host,
			"port": s.config.Gateway.Port,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
