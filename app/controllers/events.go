package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/vsokoltsov/beeg/app/channels"
	"github.com/vsokoltsov/beeg/app/utils"
)

// EventParams stores json parameters of the event
type EventParams struct {
	ID    int
	Label string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WsEndpoint provides websockets connection
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var params EventParams
		// _, message, err := c.ReadMessage()
		err := c.ReadJSON(&params)
		if err != nil {
			log.Println("read:", err)
			break
		}
		var eventID = strconv.Itoa(params.ID)
		var redisLabel = strings.Join([]string{eventID, params.Label}, utils.RedisKeySeparator)
		channels.Labels <- redisLabel
	}
}
