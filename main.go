package main

import (
	"os"

	"github.com/vsokoltsov/beeg/app"
)

func main() {
	appEnv := os.Getenv("APP_ENV")
	app := app.App{}
	app.Initialize(appEnv)
	app.Start()
}
