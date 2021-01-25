package main

import (
	"flag"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			randId := rand.Intn(10)
			randLabel := rand.Intn(15)
			strLabel := strconv.Itoa(randLabel)
			writeJsonError := c.WriteJSON(map[string]interface{}{"id": randId, "label": strLabel})
			if writeJsonError != nil {
				log.Fatal("json write: ", writeJsonError)
			}
		}
	}()

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
