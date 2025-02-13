package logger

import (
	"fmt"
	"log"
	"os"
)

func Logger(message interface{}) {
	var logMessage string
	switch v := message.(type) {
	case string:
		logMessage = v
	case error:
		logMessage = v.Error()
	default:
		logMessage = fmt.Sprintf("Unknown type: %v", v)
	}

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(logMessage)

	log.SetOutput(os.Stdout)
	log.Println(logMessage)
}
