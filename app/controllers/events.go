package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/vsokoltsov/beeg/app/utils"

	"github.com/vsokoltsov/beeg/app/channels"
)

// EventParams stores json parameters of the event
type EventParams struct {
	ID    int
	Label string
}

// CreateEvent create or update existing event
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params EventParams
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)
	if params.ID == 0 || len(params.Label) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"errors": []string{"Invalid parameters"}})
		return
	}
	var eventID = strconv.Itoa(params.ID)
	var redisLabel = strings.Join([]string{eventID, params.Label}, utils.RedisKeySeparator)
	channels.Labels <- redisLabel
	result := <-channels.RedisItems
	json.NewEncoder(w).Encode(map[string]interface{}{"id": params.ID, "label": params.Label, "viewed": result.Value})
}
