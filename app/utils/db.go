package utils

import (
	"log"
	"strings"

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

type event struct {
	id    string `db:"id"`
	label string `db:"label"`
	count int    `db:"count"`
}

func bulkSaveToDB(events []event) {
	tx := DB.MustBegin()
	stmt, prepareErr := tx.Prepare("insert into beeg_events(id, label, count) values(?, ?, ?) on duplicate key update count = count + ?")
	if prepareErr != nil {
		log.Fatal("PREPARE FOR TRANSACTION ERROR", prepareErr)
	}
	for _, event := range events {
		_, stmtErr := stmt.Exec(event.id, event.label, event.count, event.count)
		if stmtErr != nil {
			log.Fatal(event.id, event.label, event.count, stmtErr)
		}
	}
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
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

// redisScan scans all redis keys
func redisScan() []string {
	var (
		cursor uint64
		num    uint64 = 100
		keys   []string
		err    error
	)

	for {
		keys, cursor, err = RedisClient.Scan(cursor, "*-rediskeyseparator-*", int64(num)).Result()
		if err != nil {
			log.Fatal("Redis retrieve keys: ", err)
		}
		if cursor == 0 {
			break
		}
	}
	return keys
}

// getRedisData forms map with key and value
func getRedisData() map[string]int {
	var (
		data = make(map[string]int)
		keys []string
	)
	keys = redisScan()
	for _, key := range keys {
		redisVal, _ := RedisClient.Get(key).Int()
		data[key] = redisVal
	}
	return data
}

// FromCacheToDB saves existing events to database every 10 seconds
func FromCacheToDB() {

	var labels []event
	var data = getRedisData()
	for key, value := range data {
		stringKey := key
		stringKeySplit := strings.Split(stringKey, RedisKeySeparator)
		if len(stringKeySplit) < 2 {
			continue
		}
		id := stringKeySplit[0]
		label := stringKeySplit[1]
		if len(id) == 0 || len(label) == 0 {
			continue
		}
		eventParam := event{
			id:    id,
			label: label[:100],
			count: value,
		}
		labels = append(labels, eventParam)
	}
	go bulkSaveToDB(labels)
	labels = []event{}
	RedisClient.Do("FLUSHALL").Result()
}
