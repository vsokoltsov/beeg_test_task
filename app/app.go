package app

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/vsokoltsov/beeg/app/channels"
	"github.com/vsokoltsov/beeg/app/controllers"
	"github.com/vsokoltsov/beeg/app/utils"
)

const (
	developmentTimeout = 10
	testTimeout        = 0
)

type App struct {
	Router           *mux.Router
	ConnectionString string
}

// Initialize populates App struct with
// necessary parameters
func (app *App) Initialize(env string) {
	dbConName := utils.GetDatabaseConnection(env)
	app.ConnectionString = os.Getenv(dbConName)
	utils.InitDB(app.ConnectionString)
	utils.InitRedis()

	timeout := getTimeout(env)
	go channels.ManageOperations()
	go utils.FromCacheToDB(timeout)

	router := mux.NewRouter()
	router.HandleFunc("/", controllers.CreateEvent).Methods("POST")
	app.Router = router
}

// Start run application
func (app *App) Start() {
	log.Println("Starting server on port 8000")
	err := http.ListenAndServe(":8000", app.Router)
	if err != nil {
		log.Fatalln(err)
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
