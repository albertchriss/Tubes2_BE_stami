package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/albertchriss/Tubes2_BE_stami/internal/scraper"
	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan scraper.TreeNode
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan scraper.TreeNode),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (manager *ClientManager) Run() {
	log.Println("WebSocket ClientManager running")
	for {
		select {
		case conn := <-manager.register:
			manager.mu.Lock()
			manager.clients[conn] = true
			manager.mu.Unlock()
			log.Printf("Client registered via WebSocket. Total clients: %d", len(manager.clients))

		case conn := <-manager.unregister:
			manager.mu.Lock()
			if _, ok := manager.clients[conn]; ok {
				delete(manager.clients, conn)
				log.Printf("Client unregistered from WebSocket. Total clients: %d", len(manager.clients))
			}
			manager.mu.Unlock()

		case treeNodeData := <-manager.broadcast:
			manager.mu.Lock()
			if treeNodeData.Name == "" && len(treeNodeData.Children) == 0 {
				log.Println("Attempted to broadcast an empty or invalid TreeNode. Skipping.")
				manager.mu.Unlock()
				continue
			}

			jsonData, err := json.Marshal(treeNodeData)
			if err != nil {
				log.Printf("Error marshalling TreeNode to JSON for broadcast: %v", err)
				manager.mu.Unlock()
				continue
			}

			log.Printf("Broadcasting TreeNode update (root: %s) to %d clients.", treeNodeData.Name, len(manager.clients))

			disconnectedClients := []*websocket.Conn{}
			for client := range manager.clients {
				if err := client.WriteMessage(websocket.TextMessage, jsonData); err != nil {
					log.Printf("Error writing message to WebSocket client %p: %v", client, err)
					disconnectedClients = append(disconnectedClients, client)
				}
			}

			// Hapus klien yang terputus dari daftar
			for _, client := range disconnectedClients {
				if _, ok := manager.clients[client]; ok {
					delete(manager.clients, client)
					client.Close()
					log.Printf("WebSocket client %p automatically unregistered due to write error.", client)
				}
			}
			manager.mu.Unlock()
		}
	}
}

// BroadcastNode mengirimkan TreeNode ke channel broadcast.
func (manager *ClientManager) BroadcastNode(node scraper.TreeNode) {
	if node.Name == "" && len(node.Children) == 0 {
		log.Println("BroadcastNode: Received an empty or invalid node, not broadcasting.")
		return
	}
	log.Printf("Node (root: %s) received by BroadcastNode, attempting to send to broadcast channel.", node.Name)
	manager.broadcast <- node
}
