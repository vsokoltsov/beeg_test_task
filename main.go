package main

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vsokoltsov/beeg/app"
)

func main() {
	appEnv := os.Getenv("APP_ENV")
	app := app.App{}
	app.Initialize(appEnv)
	app.Start()
}
