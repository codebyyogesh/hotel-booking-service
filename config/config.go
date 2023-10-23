package env

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

//const projectDirName = "hotel-booking-service"
const projectDirName = "/app"
func LoadEnv() {
    fmt.Println("My Golang")
    projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
    currentWorkDirectory, _ := os.Getwd()
    rootPath := projectName.Find([]byte(currentWorkDirectory))
    err := godotenv.Load(string(rootPath) + `/.env`)

    if err != nil {
        log.Fatalf("Error loading .env file")
        log.Fatalln("Load env", projectName, currentWorkDirectory, rootPath)
    }
}