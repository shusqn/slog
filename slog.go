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

var currentLogLevel int8 = LOG_LEVEL_DEBUG

var errorLog, infoLog, debugLog, warnLog *log.Logger
var errorFile, infoFile, debugFile, warnFile *os.File

func SetCurrentLogLevel(lv int8) {
	currentLogLevel = lv
}

//build
func buildLog(fileName string, prefix string, oldFile *os.File) (*os.File, *log.Logger) {
	logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	//多个地方同时写入
	newlog := log.New(io.MultiWriter(logFile, os.Stderr),
		prefix,
		log.Lshortfile|log.LstdFlags|log.LUTC|log.Lmicroseconds)

	if oldFile != nil {
		oldFile.Close()
	}
	return logFile, newlog
}

//
func buildFileName(dirName string) string {
	fileName := dirName
	fileName += "/"
	fileName += time.Now().Format("20060102-15")
	fileName += ".log"

	return fileName
}

//error
func Warn(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_WARN {
		return
	}

	dirName := time.Now().Format("2006-01-02")
	fileName := buildFileName(dirName)
	if warnFile == nil || warnFile.Name() != fileName {
		os.MkdirAll(dirName, 0777)
		_logfile, _log := buildLog(fileName, "[WARN ] ", warnFile)
		warnFile = _logfile
		warnLog = _log
	}

	warnLog.Output(2, fmt.Sprintln(v...))
}

//error
func Info(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_INFO {
		return
	}
	dirName := time.Now().Format("2006-01-02")
	fileName := buildFileName(dirName)
	if warnFile == nil || warnFile.Name() != fileName {
		os.MkdirAll(dirName, 0777)
		_logfile, _log := buildLog(fileName, "[INFO ] ", infoFile)
		infoFile = _logfile
		infoLog = _log
	}

	infoLog.Output(2, fmt.Sprintln(v...))
}

//error
func Error(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_ERROR {
		return
	}
	dirName := time.Now().Format("2006-01-02")
	fileName := buildFileName(dirName)
	if warnFile == nil || warnFile.Name() != fileName {
		os.MkdirAll(dirName, 0777)
		_logfile, _log := buildLog(fileName, "[ERROR] ", errorFile)
		errorFile = _logfile
		errorLog = _log
	}

	errorLog.Output(2, fmt.Sprintln(v...))
}

//debug
func Debug(v ...interface{}) {
	if currentLogLevel < LOG_LEVEL_DEBUG {
		return
	}
	dirName := time.Now().Format("2006-01-02")
	fileName := buildFileName(dirName)
	if warnFile == nil || warnFile.Name() != fileName {
		os.MkdirAll(dirName, 0777)
		_logfile, _log := buildLog(fileName, "[DEBUG] ", debugFile)
		debugFile = _logfile
		debugLog = _log
	}

	debugLog.Output(2, fmt.Sprintln(v...))
}
