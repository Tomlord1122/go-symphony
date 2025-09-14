package server

import (
	"net/http"
	"os"

	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

func (s *Server) RegisterRoutes() http.Handler {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	// Set Gin mode based on environment
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	// Configure CORS based on environment
	var allowedOrigins []string
	frontendURL := os.Getenv("FRONT_END_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173" // Default fallback
	}

	if appEnv == "production" {
		// In production, use the configured frontend URL
		allowedOrigins = []string{frontendURL}
	} else {
		// In development, allow configured frontend URL plus common development servers
		allowedOrigins = []string{frontendURL, "http://localhost:3000", "http://localhost:8080"}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/websocket", s.websocketHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

// websocketUpgrader configures the websocket upgrader
var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		// In production, you should implement proper origin checking
		return true
	},
}

func (s *Server) websocketHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("could not upgrade to websocket: %v", err)
		return
	}
	defer conn.Close()

	// Send periodic messages to the client
	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := conn.WriteMessage(websocket.TextMessage, []byte(payload))
		if err != nil {
			log.Printf("websocket write error: %v", err)
			break
		}
		time.Sleep(time.Second * 2)
	}
}
