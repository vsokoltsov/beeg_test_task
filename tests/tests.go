package tests

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
	"github.com/vsokoltsov/beeg/app"
	"github.com/vsokoltsov/beeg/app/utils"
)

const migrationsPath = "/app/app/migrations/"

// AppInstance saves application information
var AppInstance app.App

// MakeRequest peforms http request
func MakeRequest(verb string, path string, params *strings.Reader) *httptest.ResponseRecorder {
	body := params
	if params == nil {
		body = strings.NewReader("")
	}
	req, _ := http.NewRequest(verb, path, body)
	rr := httptest.NewRecorder()
	AppInstance.Router.ServeHTTP(rr, req)
	return rr
}

// DeleteFromRedis delete particular key from redis
func DeleteFromRedis(key string) {
	utils.RedisClient.Del(key)
}

// SetKeyToRedis sets particular key and value to redis
func SetKeyToRedis(key, value string) {
	resErr := utils.RedisClient.Set(key, value, 0).Err()
	if resErr != nil {
		panic(resErr)
	}
}

// DeleteAllRedisKeys delete all keys in redis service
func DeleteAllRedisKeys() {
	keys, redisKeysErr := utils.RedisClient.Do("KEYS", "*").Result()
	if redisKeysErr != nil {
		log.Fatal(redisKeysErr)
	}
	for _, key := range keys.([]interface{}) {
		utils.RedisClient.Del(key.(string))
	}
}

// TestMain represents default wrapper for all tests
func TestMain(m *testing.M) {
	appEnv := os.Getenv("APP_ENV")

	AppInstance = app.App{}
	AppInstance.Initialize(appEnv)

	db, err := sql.Open("mysql",
		AppInstance.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	errDB := goose.Run("up", db, migrationsPath)
	if errDB != nil {
		log.Fatal("DATABASE ERROR ", errDB)
	}

	exitVal := m.Run()

	os.Exit(exitVal)

}
