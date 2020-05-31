package utils

import (
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

// DB saves db connection info to instance
var DB *sqlx.DB

const (
	developmentDBConString = "DB_CON"
	testDBString           = "DB_TEST_CON"
)

// InitDB initializes database instance
func InitDB(dataSource string) {
	var err error
	DB, err = sqlx.Connect("mysql", dataSource)
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

// GetDatabaseConnection returns name of the
// conection string based on env variable value
func GetDatabaseConnection(env string) string {
	switch env {
	case "development":
		return developmentDBConString
	case "test":
		return testDBString
	default:
		return developmentDBConString
	}
}

// FromCacheToDB saves existing events to database every 10 seconds
func FromCacheToDB(interval time.Duration) {
	for {
		time.Sleep(time.Second * interval)
		keys, redisKeysErr := RedisClient.Do("KEYS", "*").Result()
		if redisKeysErr != nil {
			log.Fatal(redisKeysErr)
		}
		for _, key := range keys.([]interface{}) {
			stringKey := key.(string)
			redisVal, _ := RedisClient.Get(stringKey).Int()
			stringKeySplit := strings.Split(stringKey, RedisKeySeparator)
			if len(stringKeySplit) < 2 {
				continue
			}
			id := stringKeySplit[0]
			label := stringKeySplit[1]
			if len(id) == 0 || len(label) == 0 {
				continue
			}
			saveToDB(id, label, redisVal)
		}
	}
}
