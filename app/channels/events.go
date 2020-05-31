package channels

import (
	"sync"

	"github.com/vsokoltsov/beeg/app/utils"
)

// RedisItem represents redis key / value for
// passing through channels
type RedisItem struct {
	Key   string
	Value int
}

var mutex = &sync.Mutex{}

// RedisItems provides communication between RedisItems structures
var RedisItems = make(chan RedisItem)

// Labels store information about labels
var Labels = make(chan string)

// ManageOperations select and update information from redis
func ManageOperations() {
	for {
		select {
		case label := <-Labels:
			currentVal, _ := utils.RedisClient.Get(label).Int()

			mutex.Lock()
			currentVal++
			mutex.Unlock()

			resErr := utils.RedisClient.Set(label, currentVal, 0).Err()
			if resErr != nil {
				panic(resErr)
			}

			RedisItems <- RedisItem{
				Key:   label,
				Value: currentVal,
			}
		}
	}
}
