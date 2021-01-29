package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/vsokoltsov/beeg/pkg/events"
)

func SetupRoutes(rdb *redis.Client) *mux.Router {
	r := mux.NewRouter()

	eventsHandler := events.NewHandler(rdb)
	r.HandleFunc("/ws", eventsHandler.WsEndpoint)
	return r
}
