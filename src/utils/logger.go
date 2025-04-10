package utils

import (
	"log"
	"os"
)

var (
	Logger  log.Logger
	logFile *os.File
)

func InitLogger(path string) error {
	var err error
	logFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	Logger.SetOutput(logFile)
	Logger.Printf("logger init finished with path: %s", path)
	return nil
}

func Log(format string, v ...any) {
	Logger.Printf(format, v...)
	log.Printf(format, v...)
}

func CloseLogger() {
	logFile.Close()
}
