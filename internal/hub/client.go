package hub

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	Receive  chan []byte
	PlayerID string
}

func NewClient(hub *Hub, conn *websocket.Conn, playerID string) *Client {
	return &Client{
		Hub:      hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Receive:  make(chan []byte, 256),
		PlayerID: playerID,
	}
}

const (
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		close(c.Receive)
		c.Conn.Close()
	}()
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"playerID": c.PlayerID,
			}).Info("ReadPump closing due to error")
			break
		}
		c.Receive <- message
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.WithFields(log.Fields{
					"error":    err,
					"playerID": c.PlayerID,
				}).Warn("Error writing text message")
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithFields(log.Fields{
					"error":    err,
					"playerID": c.PlayerID,
				}).Warn("Error writing ping message")
				return
			}
		}
	}
}

func (c *Client) SendJSON(v interface{}) error {
	if c == nil || c.Send == nil {
		return errors.New("invalid client")
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	select {
	case c.Send <- b:
		return nil
	default:
		return errors.New("send buffer full")
	}
}

func (c *Client) SendError(err error) error { //TODO: 引数の修正
	if c == nil || err == nil {
		return nil
	}
	payload := map[string]interface{}{
		"type":    "error",
		"message": err.Error(),
		"ts":      time.Now().Unix(),
	}
	return c.SendJSON(payload)
}
