package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"pet-social-game/server/internal/model"
	"pet-social-game/server/internal/pkg/auth"
	"pet-social-game/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ClientMsg struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	ToID    uint64 `json:"toId"`
	Content string `json:"content"`
}

type ServerMsg map[string]interface{}

type Client struct {
	UserID uint64
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

type Hub struct {
	mu      sync.RWMutex
	clients map[uint64]map[*Client]struct{}
	rdb     *redis.Client
	svc     *service.Services
	secret  string
}

func NewHub(rdb *redis.Client, svc *service.Services, secret string) *Hub {
	return &Hub{
		clients: make(map[uint64]map[*Client]struct{}),
		rdb:     rdb,
		svc:     svc,
		secret:  secret,
	}
}

func (h *Hub) Run(ctx context.Context) {
	pubsub := h.rdb.Subscribe(ctx, "pubsub:chat:world", "pubsub:chat:dm")
	ch := pubsub.Channel()
	go func() {
		<-ctx.Done()
		_ = pubsub.Close()
	}()
	for msg := range ch {
		h.dispatchPubSub(msg.Channel, []byte(msg.Payload))
	}
}

func (h *Hub) dispatchPubSub(channel string, payload []byte) {
	var view service.ChatMsgView
	if err := json.Unmarshal(payload, &view); err != nil {
		return
	}
	out, _ := json.Marshal(ServerMsg{
		"type": "chat.msg", "channel": view.Channel, "fromId": view.FromID,
		"fromName": view.FromName, "toId": view.ToID, "content": view.Content, "ts": view.Ts, "id": view.ID,
	})
	if channel == "pubsub:chat:world" || view.Channel == model.ChannelWorld {
		h.broadcastAll(out)
		return
	}
	h.sendToUser(view.FromID, out)
	if view.ToID != 0 {
		h.sendToUser(view.ToID, out)
	}
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.UserID] == nil {
		h.clients[c.UserID] = make(map[*Client]struct{})
	}
	h.clients[c.UserID][c] = struct{}{}
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.UserID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, c.UserID)
		}
	}
	close(c.Send)
}

func (h *Hub) broadcastAll(payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, set := range h.clients {
		for c := range set {
			select {
			case c.Send <- payload:
			default:
			}
		}
	}
}

func (h *Hub) sendToUser(userID uint64, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients[userID] {
		select {
		case c.Send <- payload:
		default:
		}
	}
}

func (h *Hub) Handle(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "missing token"})
		return
	}
	claims, err := auth.ParseToken(h.secret, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid token"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &Client{
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 64),
		Hub:    h,
	}
	h.register(client)
	_ = h.svc.Store.SetOnline(c.Request.Context(), claims.UserID, 2*time.Minute)

	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	ticker := time.NewTicker(25 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister(c)
		_ = c.Conn.Close()
	}()
	c.Conn.SetReadLimit(4096)
	_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		_ = c.Hub.svc.Store.SetOnline(context.Background(), c.UserID, 2*time.Minute)
		return nil
	})
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		var msg ClientMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			c.push(ServerMsg{"type": "error", "msg": "bad json"})
			continue
		}
		c.handle(msg)
	}
}

func (c *Client) push(m ServerMsg) {
	b, _ := json.Marshal(m)
	select {
	case c.Send <- b:
	default:
	}
}

func (c *Client) handle(msg ClientMsg) {
	switch msg.Type {
	case "ping":
		c.push(ServerMsg{"type": "pong"})
	case "chat.send":
		view, reason, err := c.Hub.svc.SendChat(context.Background(), c.UserID, service.SendChatInput{
			Channel: msg.Channel, ToID: msg.ToID, Content: msg.Content,
		})
		if err != nil {
			c.push(ServerMsg{"type": "error", "msg": err.Error()})
			return
		}
		if reason != "" {
			c.push(ServerMsg{"type": "chat.rejected", "reason": reason})
			return
		}
		payload, _ := json.Marshal(view)
		channel := "pubsub:chat:world"
		if view.Channel == model.ChannelFriend {
			channel = "pubsub:chat:dm"
		}
		if err := c.Hub.rdb.Publish(context.Background(), channel, payload).Err(); err != nil {
			log.Printf("pubsub publish: %v", err)
			// Fallback local deliver
			out, _ := json.Marshal(ServerMsg{
				"type": "chat.msg", "channel": view.Channel, "fromId": view.FromID,
				"fromName": view.FromName, "toId": view.ToID, "content": view.Content, "ts": view.Ts, "id": view.ID,
			})
			if view.Channel == model.ChannelWorld {
				c.Hub.broadcastAll(out)
			} else {
				c.Hub.sendToUser(view.FromID, out)
				c.Hub.sendToUser(view.ToID, out)
			}
		}
	default:
		c.push(ServerMsg{"type": "error", "msg": "unknown type"})
	}
}
