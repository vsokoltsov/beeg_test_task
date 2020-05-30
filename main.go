package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type RedisItem struct {
	Key   string
	Value int
}

var DB *sqlx.DB
var redisClient *redis.Client
var requestCount int64
var mutex = &sync.Mutex{}
var redisItems = make(chan RedisItem)
var labels = make(chan string)

type EventParams struct {
	ID    int
	Label string
}

type BeegEvent struct {
	ID    int    `db:"id" json:"id"`
	Label string `db:"label" json:"label"`
	Count int    `db:"count" json:"count"`
}

func manageOperations() {
	for {
		select {
		case label := <-labels:
			currentVal, _ := redisClient.Get(label).Int()

			mutex.Lock()
			currentVal++
			mutex.Unlock()

			resErr := redisClient.Set(label, currentVal, 0).Err()
			if resErr != nil {
				panic(resErr)
			}

			redisItems <- RedisItem{
				Key:   label,
				Value: currentVal,
			}
		}
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestCount++
	var params EventParams
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&params)

	var eventID = strconv.Itoa(params.ID)
	var redisLabel = strings.Join([]string{eventID, params.Label}, "-")
	labels <- redisLabel
	result := <-redisItems
	json.NewEncoder(w).Encode(map[string]interface{}{"id": params.ID, "label": params.Label, "viewed": result.Value})
}

func saveToDB(id string, label string, counter int) {
	tx := DB.MustBegin()
	tx.MustExec(
		"insert into beeg_events(id, label, count) values(?, ?, ?) on duplicate key update count = ?",
		id,
		label,
		counter,
		counter,
	)
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func fromCacheToDB() {
	for {
		time.Sleep(time.Second * 10)
		keys, redisKeysErr := redisClient.Do("KEYS", "*").Result()
		if redisKeysErr != nil {
			log.Fatal(redisKeysErr)
		}
		for _, key := range keys.([]interface{}) {
			stringKey := key.(string)
			redisVal, _ := redisClient.Get(stringKey).Int()
			stringKeySplit := strings.Split(stringKey, "-")
			id := stringKeySplit[0]
			label := stringKeySplit[1]
			if len(id) == 0 || len(label) == 0 {
				continue
			}
			go saveToDB(id, label, redisVal)
		}
	}
}

func main() {
	connString := os.Getenv("DB_CON")
	var err error
	DB, err = sqlx.Connect("mysql", connString)
	if err != nil {
		log.Panic(err)
	}

	if err = DB.Ping(); err != nil {
		log.Panic(err)
	}
	DB.SetMaxOpenConns(1000)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, redisErr := redisClient.Ping().Result()
	if redisErr != nil {
		fmt.Println(redisErr)
	}
	go manageOperations()
	go fromCacheToDB()
	r := mux.NewRouter()
	r.HandleFunc("/", defaultHandler).Methods("POST")
	serverError := http.ListenAndServe(":8000", r)
	if serverError != nil {
		log.Fatal(serverError)
	}
}
