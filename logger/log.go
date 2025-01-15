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

func getPrefix(level int) string {
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
	return prefix
}

func logToStdOut(msg string) {
	log.SetFlags(0)
	log.Println(msg)
}

func Log(level int, v interface{}) {
	var finalMsg string

	prefix := getPrefix(level)
	prefix = fmt.Sprintf("%s %s: %s", prefix, time.Now().Format("2006-01-02 15:04:05"), Reset)

	switch msg := v.(type) {
	case string:
		finalMsg = prefix + msg
	case error:
		finalMsg = prefix + msg.Error()
	default:
		finalMsg = prefix + fmt.Sprintf("%v", msg)
	}

	logToStdOut(finalMsg)
	if level == CRITICAL {
		os.Exit(1)
	}
}

func Logf(level int, format string, v ...interface{}) {
	var finalMsg string
	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}

	prefix := getPrefix(level)
	prefix = fmt.Sprintf("%s %s: %s", prefix, time.Now().Format("2006-01-02 15:04:05"), Reset)

	finalMsg = prefix + format
	logToStdOut(finalMsg)

	if level == CRITICAL {
		os.Exit(1)
	}
}

func Info(msg string) {
	prefix := fmt.Sprintf("%s[INFO] %s: %s", Green, time.Now().Format("2006-01-02 15:04:05"), Reset)
	finalMsg := prefix + msg
	logToStdOut(finalMsg)
}

func InfoF(msg string, v ...interface{}) {
	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}
	prefix := fmt.Sprintf("%s[INFO] %s: %s", Green, time.Now().Format("2006-01-02 15:04:05"), Reset)
	finalMsg := prefix + msg
	logToStdOut(finalMsg)
}

func Critical(msg string) {
	prefix := fmt.Sprintf("%s[CRITICAL] %s", Red, time.Now().Format("2006-01-02 15:04:05"))
	finalMsg := prefix + msg
	logToStdOut(finalMsg)
	os.Exit(1)
}

func CriticalF(msg string, v ...interface{}) {
	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}
	prefix := fmt.Sprintf("%s[CRITICAL] %s", Red, time.Now().Format("2006-01-02 15:04:05"))
	finalMsg := prefix + msg
	logToStdOut(finalMsg)
	os.Exit(1)
}
