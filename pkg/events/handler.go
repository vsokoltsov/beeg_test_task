package events

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type eventParams struct {
	id    int    `json:"id"`
	label string `json:"label"`
}

type Handler struct {
	redisClient *redis.Client
}

func NewHandler(rd *redis.Client) Handler {
	return Handler{
		redisClient: rd,
	}
}

func (h *Handler) WsEndpoint(w http.ResponseWriter, r *http.Request) {

}
