package logwrapper

import (
	"fmt"
	"log"
	"time"
)

type LogLevel int8

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	FATAL
)

var logLevelStr = []string{"DEBUG", "INFO", "ERROR", "FATAL"}

func (l LogLevel) String() string {
	if (l < DEBUG) || (l > FATAL) {
		return "???"
	} else {
		return logLevelStr[l]
	}
}

type LogObserver interface {
	Log(logLevel LogLevel, msg string)
}

func SetLogLevel(newLogLevel LogLevel) {
	globalLogLevel = newLogLevel
}

func SetObserver(observer *LogObserver, allLevel bool) {
	logObserver = observer
	observeAllLevel = allLevel
}

func LogTimestamps(useTimestamps bool) {
	logTimestamps = useTimestamps
}

var globalLogLevel LogLevel
var logObserver *LogObserver
var observeAllLevel bool
var logTimestamps bool

func createLogMsg(logLevel LogLevel, callingContext string, msg string, err *error) string {
	logMsg := msg
	if err != nil {
		logMsg = fmt.Sprintf("%s (%v)", msg, err)
	}
	if logTimestamps {
		timestamp := time.Now().Format(time.RFC3339)
		return fmt.Sprintf("%s [%s] %s - %s", timestamp, callingContext, logLevel.String(), logMsg)
	} else {
		return fmt.Sprintf("[%s] %s - %s", callingContext, logLevel.String(), logMsg)
	}
}

func implLogging(logLevel LogLevel, callingContext string, msg string, err *error) {
	matchesGlobalLevel := globalLogLevel >= logLevel
	var logMsg string
	if matchesGlobalLevel {
		logMsg = createLogMsg(logLevel, callingContext, msg, err)
		log.Println(logMsg)
	}
	if (logObserver != nil) && (observeAllLevel || matchesGlobalLevel) {
		if logMsg == "" {
			logMsg = createLogMsg(logLevel, callingContext, msg, err)
		}
		(*logObserver).Log(logLevel, logMsg)
	}
	if logLevel == FATAL {
		panic(logMsg)
	}
}

func Debug(callingContext string, msg string) {
	implLogging(DEBUG, callingContext, msg, nil)
}

func DebugErr(callingContext string, msg string, err error) {
	implLogging(DEBUG, callingContext, msg, &err)
}

func Info(callingContext string, msg string) {
	implLogging(INFO, callingContext, msg, nil)
}

func InfoErr(callingContext string, msg string, err error) {
	implLogging(INFO, callingContext, msg, &err)
}

func Error(callingContext string, msg string) {
	implLogging(ERROR, callingContext, msg, nil)
}

func ErrorErr(callingContext string, msg string, err error) {
	implLogging(ERROR, callingContext, msg, &err)
}

func Fatal(callingContext string, msg string) {
	implLogging(FATAL, callingContext, msg, nil)
}

func FatalErr(callingContext string, msg string, err error) {
	implLogging(FATAL, callingContext, msg, &err)
}
