package middleware

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger() {
	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogInfo(message string) {
	Logger.Println("INFO: " + message)
}

func LogError(message string, err error) {
	Logger.Println("ERROR: " + message, err)
}