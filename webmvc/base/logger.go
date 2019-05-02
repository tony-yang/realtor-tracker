package base

import (
	"log"
	"os"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

// Initial default level
var level = LevelDebug
var MVCLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

func Debug(v ...interface{}) {
	if level <= LevelDebug {
		MVCLogger.Printf("[Debug] %v\n", v)
	}
}

func Info(v ...interface{}) {
	if level <= LevelInfo {
		MVCLogger.Printf("[Info] %v\n", v)
	}
}

func Warn(v ...interface{}) {
	if level <= LevelWarning {
		MVCLogger.Printf("[Warn] %v\n", v)
	}
}

func Error(v ...interface{}) {
	if level <= LevelError {
		MVCLogger.Printf("[Error] %v\n", v)
	}
}

func Critical(v ...interface{}) {
	if level <= LevelCritical {
		MVCLogger.Print("[Critical] %v\n", v)
	}
}
