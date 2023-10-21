package env

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "hotel-booking-service"
func LoadEnv() {
    projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
    currentWorkDirectory, _ := os.Getwd()
    rootPath := projectName.Find([]byte(currentWorkDirectory))
    err := godotenv.Load(string(rootPath) + `/.env`)

    if err != nil {
        log.Fatalf("Error loading .env file")
    }
}