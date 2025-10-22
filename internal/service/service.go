package service

/*import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/shii-park/Metasugo-Backend/internal/hub"
)

/******************************************************************/
/*type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *hub
	UserID string
}

type Hub struct {
	Clients     map[*Client]bool
	register    chan *Client
	Unregister  chan *Client
	Broadcast   chan []byte
	UserClients map[string]*Client
	mu          sync.RWMutex
}

/******************************************************************/
/*func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan []byte),
		UserClients: make(map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.UserClients[client.UserID] = client
			h.mu.Unlock()
			log.Printf("Client registerd: %s (UserID: %s)", client.ID, client.UserID)
		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.UserClients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			/*case message:=<-h.Broadcast:
		}
	}
}
*/
