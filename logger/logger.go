package logger

import (
	"encoding/json"
	"io"
	"log"
)

type CustomLogger struct {
	Writer io.Writer
}

func New(writer io.Writer) *CustomLogger {
	logger := &CustomLogger{
		Writer: writer,
	}

	log.SetFlags(0)
	log.SetOutput(logger.Writer)

	return logger
}

func (logger *CustomLogger) Message(message interface{}){
	logData, _ := json.Marshal(message)

	log.Println(string(logData))
}

func (logger *CustomLogger) Fatal(message interface{}){
	logData, _ := json.Marshal(message)

	log.Fatalln(string(logData))
}

func (logger *CustomLogger) Panic(message interface{}){
	logData, _ := json.Marshal(message)

	log.Panicln(string(logData))
}
