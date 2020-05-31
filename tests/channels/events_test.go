package testchannels

import (
	"testing"

	"github.com/vsokoltsov/beeg/app/channels"
	"github.com/vsokoltsov/beeg/app/utils"
	"github.com/vsokoltsov/beeg/tests"
)

// TestMain rewrite base test for events channels module
func TestMain(m *testing.M) {
	tests.TestMain(m)
}

// TestSuccessRedisReadAndWrite success read and write operation
func TestSuccessRedisReadAndWrite(t *testing.T) {
	go channels.ManageOperations()
	key := "test_label"
	defer tests.DeleteFromRedis(key)
	channels.Labels <- key
	result := <-channels.RedisItems
	if result.Value != 1 {
		t.Error("Redis has wrong value")
	}
	utils.RedisClient.Del("")
}

// TestSuccessSavingOfExistingKey success update of existing key
func TestSuccessSavingOfExistingKey(t *testing.T) {
	go channels.ManageOperations()
	var (
		initialValue = 1
		key          = "test_label"
	)
	defer tests.DeleteFromRedis(key)

	resErr := utils.RedisClient.Set(key, initialValue, 0).Err()
	if resErr != nil {
		panic(resErr)
	}
	channels.Labels <- key
	result := <-channels.RedisItems
	initialValue++
	if result.Value != initialValue {
		t.Error("Existing redis key has wrong value")
	}
}
