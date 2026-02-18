package logger

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
}

var instance *Logger

func Init(prefix string) {
	instance = &Logger{
		info:  log.New(os.Stdout, "[INFO] "+prefix+" ", log.LstdFlags|log.Lshortfile),
		error: log.New(os.Stderr, "[ERROR] "+prefix+" ", log.LstdFlags|log.Lshortfile),
		debug: log.New(os.Stdout, "[DEBUG] "+prefix+" ", log.LstdFlags|log.Lshortfile),
	}
}

func Info(v ...interface{}) {
	if instance != nil {
		instance.info.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if instance != nil {
		instance.info.Printf(format, v...)
	}
}

func Error(v ...interface{}) {
	if instance != nil {
		instance.error.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if instance != nil {
		instance.error.Printf(format, v...)
	}
}

func Debug(v ...interface{}) {
	if instance != nil {
		instance.debug.Println(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if instance != nil {
		instance.debug.Printf(format, v...)
	}
}
