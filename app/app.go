package app

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vsokoltsov/beeg/app/channels"
	"github.com/vsokoltsov/beeg/app/controllers"
	"github.com/vsokoltsov/beeg/app/utils"
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

	go channels.ManageOperations()
	go utils.FromCacheToDB()

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
