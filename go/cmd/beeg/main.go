package main

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
	"github.com/vsokoltsov/beeg/pkg/app"
)

func getEnv(key, def string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return def
	}
	return value
}

func getProjectPath(delimeter string) string {
	projectDirectory, directoryErr := os.Getwd()

	if directoryErr != nil {
		log.Fatalf("Could not locate current directory: %s", directoryErr)
	}

	isUnderCmd := strings.Contains(projectDirectory, delimeter)
	if isUnderCmd {
		var cmdIdx int
		splitPath := strings.Split(projectDirectory, "/")
		for idx, pathElem := range splitPath {
			if pathElem == delimeter {
				cmdIdx = idx
				break
			}
		}
		projectDirectory = strings.Join(splitPath[:cmdIdx], "/")
	}

	return projectDirectory
}

func main() {
	pathDelimiter := getEnv("PATH_SEPARATOR", "cmd")
	projectPath := getProjectPath(pathDelimiter)
	err := godotenv.Load(path.Join(projectPath, ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		env           = getEnv("APP_ENV", "development")
		appPort       = getEnv("APP_PORT", "8080")
		appHost       = getEnv("APP_HOST", "localhost")
		sqlDbUser     = getEnv("MYSQL_USER", "user")
		sqlDbPassword = getEnv("MYSQL_PASSWORD", "password")
		sqlDbHost     = getEnv("MYSQL_HOST", "localhost")
		sqlDbPort     = getEnv("MYSQL_PORT", "3306")
		sqlDbName     = getEnv("MYSQL_DATABASE", "beeg_go")
		dbProvider    = getEnv("DB_PROVIDER", "mysql")
		redisHost     = getEnv("REDIS_HOST", "localhost")
		redisPort     = getEnv("REDIS_PORT", "6379")
		redisPassword = getEnv("REDIS_PASSWORD", "")
	)

	dbConnUser := sqlDbUser + ":" + sqlDbPassword
	dbConnURL := sqlDbHost + ":" + sqlDbPort

	dbConnString := dbConnUser + "@tcp(" + dbConnURL + ")/" + sqlDbName + "?"

	redisConnectString := strings.Join(
		[]string{redisHost, redisPort}, ":",
	)

	app := app.App{}
	app.Initialize(env, appHost, appPort, pathDelimiter, dbProvider, dbConnString, redisConnectString, redisPassword)
	app.Run()
}
