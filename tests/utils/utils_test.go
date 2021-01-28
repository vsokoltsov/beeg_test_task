package utilstest

import (
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vsokoltsov/beeg/app/utils"
	"github.com/vsokoltsov/beeg/tests"
)

// TestMain rewrite base test for events controller module
func TestMain(m *testing.M) {
	tests.TestMain(m)
}

// BeegEvent represents beeg event model
type BeegEvent struct {
	id    int    `db:"id"`
	label string `db:"label"`
	count int    `db:"count"`
}

// flushDB Remove all items from beeg_events table
func flushDB() {
	tx := utils.DB.MustBegin()
	tx.MustExec(
		"delete from beeg_events",
	)
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// saveEventToDB saves or updates particular event to database
func saveEventToDB(id int, label string) {

	tx := utils.DB.MustBegin()
	tx.MustExec(
		"insert into beeg_events(id, label) values(?, ?) on duplicate key update count = count + 1",
		id,
		label,
	)
	// fmt.Println("EVENTS ARE ", events)
	err := tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// getEventsCount return number of events at current time
func getEventsCount() int {
	var eventsCount int
	rows := utils.DB.QueryRow("select count(*) from beeg_events")
	err := rows.Scan(&eventsCount)
	if err != nil {
		log.Fatal(err)
	}
	return eventsCount
}

// TestSuccessSavingFromRedisToDB test success saving new events to db
func TestSuccessSavingFromRedisToDB(t *testing.T) {
	flushDB()
	tests.DeleteAllRedisKeys()

	defer flushDB()
	defer tests.DeleteAllRedisKeys()

	keys := []string{
		strings.Join([]string{"1", "labeltest"}, utils.RedisKeySeparator),
		strings.Join([]string{"1", "labeltes2"}, utils.RedisKeySeparator),
	}
	for idx, k := range keys {
		tests.SetKeyToRedis(k, strconv.Itoa(1+idx))
	}
	currentEventsCount := getEventsCount()
	go utils.FromCacheToDB(0)
	time.Sleep(time.Second * 10)
	newEventsCount := getEventsCount()

	if newEventsCount != currentEventsCount+2 {
		t.Error("New events from redis were not saved to database")
	}
}

// TestSuccessSavingFromRedisToDB test success
// updating of existing events in db
func TestSuccessUpdateFromRedisToDB(t *testing.T) {
	flushDB()
	tests.DeleteAllRedisKeys()

	defer flushDB()
	defer tests.DeleteAllRedisKeys()

	lbl := strings.Join([]string{"1", "labeltest"}, utils.RedisKeySeparator)
	saveEventToDB(1, "labeltest")

	prevEventCount := 1

	keys := []string{
		lbl,
		strings.Join([]string{"1", "labeltes2"}, utils.RedisKeySeparator),
	}

	for idx, k := range keys {
		tests.SetKeyToRedis(k, strconv.Itoa(1+idx))
	}
	prevEventsCount := getEventsCount()
	go utils.FromCacheToDB(0)
	time.Sleep(time.Second * 10)
	newEventsCount := getEventsCount()

	var eventCount int
	err := utils.DB.Get(&eventCount, "select count from beeg_events where id = ? and label = ?", "1", "labeltest")
	if err != nil {
		log.Fatal(err)
	}

	if newEventsCount != prevEventsCount+1 {
		t.Error("New events from redis were not saved to database")
	}
	if eventCount != prevEventCount {
		t.Error("Number of events does not match")
	}
}
