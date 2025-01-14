package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Green   = "\033[32m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
)

func Log(level int, v interface{}) {
	var prefix string

	switch level {
	case DEBUG:
		prefix = Blue + "[DEBUG]"
	case INFO:
		prefix = Green + "[INFO]"
	case WARNING:
		prefix = Yellow + "[WARNING]"
	case ERROR:
		prefix = Red + "[ERROR]"
	case CRITICAL:
		prefix = Magenta + "[CRITICAL]"
	default:
		prefix = "[UNKNOWN]"
	}

	prefix = fmt.Sprintf("%s %s: %s", prefix, time.Now().Format("2006-01-02 15:04:05"), Reset)
	log.SetFlags(0)

	switch msg := v.(type) {
	case string:
		log.Println(prefix + msg)
	case error:
		log.Println(prefix + msg.Error())
	default:
		log.Println(prefix + fmt.Sprint(msg))
	}
	if level == CRITICAL {
		os.Exit(1)
	}
}
