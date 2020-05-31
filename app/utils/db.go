package utils

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// DB saves db connection info to instance
var DB *sqlx.DB

// InitDB initializes database instance
func InitDB() {
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

// FromCacheToDB saves existing events to database every 10 seconds
func FromCacheToDB() {
	for {
		time.Sleep(time.Second * 10)
		keys, redisKeysErr := RedisClient.Do("KEYS", "*").Result()
		if redisKeysErr != nil {
			log.Fatal(redisKeysErr)
		}
		for _, key := range keys.([]interface{}) {
			stringKey := key.(string)
			redisVal, _ := RedisClient.Get(stringKey).Int()
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
