package main

import (
	"log"
	"net/http"

	"github.com/vsokoltsov/beeg/app/controllers"

	"github.com/vsokoltsov/beeg/app/channels"

	"github.com/vsokoltsov/beeg/app/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	utils.InitDB()
	utils.InitRedis()

	go channels.ManageOperations()
	go utils.FromCacheToDB()
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.CreateEvent).Methods("POST")
	serverError := http.ListenAndServe(":8000", r)
	if serverError != nil {
		log.Fatal(serverError)
	}
}
