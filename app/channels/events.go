package channels

import (
	"os"
	"time"

	"github.com/vsokoltsov/beeg/app/utils"
)

const (
	developmentTimeout = 10
	testTimeout        = 0
)

// RedisItem represents redis key / value for
// passing through channels
type RedisItem struct {
	Key   string
	Value int
}

// RedisItems provides communication between RedisItems structures
var RedisItems = make(chan RedisItem)

// Labels store information about labels
var Labels = make(chan string)

// LabelIncremented store information about labels that were incremented
var LabelIncremented = make(chan bool)
var env = os.Getenv("APP_ENV")
var dateTicker = time.NewTicker(getTimeout(env) * time.Second)

// ManageOperations select and update information from redis
func ManageOperations() {
	for {
		select {
		case label := <-Labels:
			utils.RedisClient.Incr(label).Val()
			LabelIncremented <- true
		case <-dateTicker.C:
			utils.FromCacheToDB()
		}
	}
}

func getTimeout(env string) time.Duration {
	switch env {
	case "development":
		return developmentTimeout
	case "test":
		return testTimeout
	default:
		return developmentTimeout
	}
}
