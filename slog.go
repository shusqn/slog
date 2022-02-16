package slog

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	LOG_LEVEL_ERROR = 1
	LOG_LEVEL_WARN  = 2
	LOG_LEVEL_INFO  = 3
	LOG_LEVEL_DEBUG = 4
)

type Slogger struct {
	logger  *log.Logger
	logfile *os.File
	expire  int64
}

var currentLogLevel int8 = LOG_LEVEL_DEBUG
var farmatStr = "20060102-15"

var errorLogger, infoLogger, debugLogger, warnLogger *Slogger

func SetCurrentLogLevel(lv int8) {
	currentLogLevel = lv
}

//buildLog
func buildLog(prefix string, logger *Slogger) *Slogger {
	dirName := time.Now().Format("2006-01-02")
	fileName := buildFileName(dirName)
	os.MkdirAll(dirName, 0777)
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	//多个地方同时写入
	newlog := log.New(io.MultiWriter(logFile, os.Stderr),
		prefix,
		log.Lshortfile|log.LstdFlags|log.LUTC|log.Lmicroseconds)

	if logger != nil {
		logger.logfile.Close()
	}

	hourtimeStr := time.Now().Format(farmatStr)
	millstime, _ := time.ParseInLocation(farmatStr, hourtimeStr, time.Local)
	expire := millstime.UnixMilli() + int64(3600*1000)
	return &Slogger{newlog, logFile, expire}
}

//buildFileName
func buildFileName(dirName string) string {
	fileName := dirName
	fileName += "/"
	fileName += time.Now().Format(farmatStr)
	fileName += ".log"

	return fileName
}

//warn
func warn(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_WARN {
		return
	}

	if warnLogger == nil || warnLogger.expire <= time.Now().UnixMilli() {
		logger := buildLog("[WARN ] ", warnLogger)
		warnLogger = logger
	}
	warnLogger.logger.Output(3, fmt.Sprintln(v...))
}

//info
func info(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_INFO {
		return
	}

	if infoLogger == nil || infoLogger.expire <= time.Now().UnixMilli() {
		logger := buildLog("[INFO ] ", infoLogger)
		infoLogger = logger
	}
	infoLogger.logger.Output(3, fmt.Sprintln(v...))
}

//error
func error(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_ERROR {
		return
	}
	if errorLogger == nil || errorLogger.expire <= time.Now().UnixMilli() {
		logger := buildLog("[ERROR] ", infoLogger)
		errorLogger = logger
	}
	errorLogger.logger.Output(3, fmt.Sprintln(v...))
}

//debug
func debug(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_DEBUG {
		return
	}
	if debugLogger == nil || debugLogger.expire <= time.Now().UnixMilli() {
		logger := buildLog("[DEBUG] ", infoLogger)
		errorLogger = logger
	}
	debugLogger.logger.Output(3, fmt.Sprintln(v...))
}

//Error
func Error(v ...interface{}) {
	error(v...)
}

//Infof
func Info(v ...interface{}) {
	info(v...)
}

//Debugf
func Debug(v ...interface{}) {
	debug(v...)
}

//Warnf
func Warn(v ...interface{}) {
	warn(v...)
}

//Errorf
func Errorf(format string, v ...interface{}) {
	error(fmt.Sprintf(format, v...))
}

//Infof
func Infof(format string, v ...interface{}) {
	info(fmt.Sprintf(format, v...))
}

//Debugf
func Debugf(format string, v ...interface{}) {
	debug(fmt.Sprintf(format, v...))
}

//Warnf
func Warnf(format string, v ...interface{}) {
	warn(fmt.Sprintf(format, v...))
}
