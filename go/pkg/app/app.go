package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vsokoltsov/beeg/pkg/api"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
)

var upgrader = websocket.Upgrader{}

type App struct {
	host          string
	port          string
	env           string
	pathDelimiter string
	router        *mux.Router
	server        *http.Server
}

// Initialize populates App struct with
// necessary parameters
func (app *App) Initialize(env, host, port, pathDelimiter, dbProvider, sqlDbConnStr, redisConnectString, redisPassword string) {
	var (
		sqlDB        *sqlx.DB
		rdb          *redis.Client
		sqlDBConnErr error
	)
	app.env = env
	app.host = host
	app.port = port
	app.pathDelimiter = pathDelimiter

	url := strings.Join([]string{app.host, app.port}, ":")

	if dbProvider == "mysql" {
		sqlDbConnStr += "&charset=utf8&interpolateParams=true"
		sqlDB, sqlDBConnErr = sqlx.Connect(dbProvider, sqlDbConnStr)
		if sqlDBConnErr != nil {
			fmt.Println(sqlDbConnStr)
			log.Fatalf("Error sql database open: %s", sqlDBConnErr)
			return
		}
		sqlDB.SetMaxOpenConns(10)
		pingErr := sqlDB.Ping()
		if pingErr != nil {
			log.Fatalf("Error sql database connection: %s", pingErr)
		}
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisConnectString,
		Password: redisPassword,
		DB:       0,
	})
	app.router = api.SetupRoutes(rdb)

	app.server = &http.Server{
		Handler:      handlers.LoggingHandler(os.Stdout, app.router),
		Addr:         url,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// dbConName := utils.GetDatabaseConnection(env)
	// app.ConnectionString = os.Getenv(dbConName)
	// utils.InitDB(app.ConnectionString)
	// utils.InitRedis()

	// go channels.ManageOperations()

	// router := mux.NewRouter()
	// router.HandleFunc("/echo", controllers.WsEndpoint)
	// app.Router = router
}

// Start run application
func (app *App) Run() {
	log.Printf("Starting web server on port %s...", app.port)
	log.Fatal(app.server.ListenAndServe())
	// err := http.ListenAndServe(":8000", app.Router)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}
