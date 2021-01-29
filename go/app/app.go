package app

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vsokoltsov/beeg/app/channels"
	"github.com/vsokoltsov/beeg/app/controllers"
	"github.com/vsokoltsov/beeg/app/utils"
)

var upgrader = websocket.Upgrader{}

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

	go channels.ManageOperations()

	router := mux.NewRouter()
	router.HandleFunc("/echo", controllers.WsEndpoint)
	app.Router = router
}

// Start run application
func (app *App) Start() {

	err := http.ListenAndServe(":8000", app.Router)
	if err != nil {
		log.Fatalln(err)
	}
}
