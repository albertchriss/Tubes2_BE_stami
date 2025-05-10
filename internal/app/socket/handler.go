package socket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader untuk WebSocket
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Di production, validasi origin
		// return origin == "http://localhost:3000"
		log.Printf("WebSocket CheckOrigin: Allowing origin %s", r.Header.Get("Origin"))
		return true
	},
}

type Handler struct {
	wsManager *ClientManager
}

func NewHandler(wsManager *ClientManager) *Handler {
	return &Handler{
		wsManager: wsManager,
	}
}

// WebSocketConnectHandler menangani koneksi WebSocket baru
func (h *Handler) WebSocketConnectHandler(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade HTTP to WebSocket: %v", err)
		return
	}

	h.wsManager.register <- conn
	log.Printf("Client %p connected via WebSocket", conn.RemoteAddr())
	// Jalankan goroutine untuk membaca pesan dari koneksi WebSocket
	go h.readPump(conn)
}

// readPump membaca pesan dari koneksi WebSocket (jika ada) dan mendeteksi penutupan.
func (h *Handler) readPump(conn *websocket.Conn) {

	// Pastikan klien di-unregister dan koneksi ditutup saat fungsi ini berakhir
	defer func() {
		h.wsManager.unregister <- conn
		conn.Close()
		log.Printf("Client %p disconnected (readPump exited)", conn.RemoteAddr())
	}()

	// Atur batas baca atau properti koneksi lainnya jika perlu
	// conn.SetReadLimit(maxMessageSize)
	// conn.SetReadDeadline(time.Now().Add(pongWait))
	// conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// loop untuk deteksi penutupan koneksi
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket client %p unexpected close error: %v", conn.RemoteAddr(), err)
			} else {
				// Ini bisa jadi penutupan normal atau error lain
				log.Printf("WebSocket client %p read error (likely closed): %v", conn.RemoteAddr(), err)
			}
			break // Keluar dari loop jika ada error atau koneksi ditutup
		}
	}
}
