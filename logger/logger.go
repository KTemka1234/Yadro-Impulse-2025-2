package logger

import (
	"io"
	"log"
	"os"
)

var EventsLog = map[int]string{
	1:  "[%s] The competitor(%d) registered",
	2:  "[%s] The start time for the competitor(%d) was set by a draw to %s",
	3:  "[%s] The competitor(%d) is on the start line",
	4:  "[%s] The competitor(%d) has started",
	5:  "[%s] The competitor(%d) is on the firing range(%d)",
	6:  "[%s] The target(%d) has been hit by competitor(%d)",
	7:  "[%s] The competitor(%d) left the firing range",
	8:  "[%s] The competitor(%d) entered the penalty laps",
	9:  "[%s] The competitor(%d) left the penalty laps",
	10: "[%s] The competitor(%d) ended the main lap",
	11: "[%s] The competitor(%d) can`t continue: %s",
	32: "[%s] The competitor(%d) is disqualified",
	33: "[%s] The competitor(%d) has finished",
}

var logFile *os.File

func InitLogger(filename string) *log.Logger {
	// see chmod 0644 for more info (https://chmodcommand.com/chmod-0644/)
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger := log.Default()
		logger.Printf("Failed to open file to write logs. Only console log available: %v", err.Error())
		return logger
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", 0)
	return logger
}

func Cleanup() {
	if logFile != nil {
		logFile.Close()
	}
}
