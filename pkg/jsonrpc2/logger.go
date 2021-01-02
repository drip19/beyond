package rpc

import (
	"log"
	"os"
)

var IsDebug = false

// Logger is implemented by any logging system that is used for standard logs.
type Logger interface {
	Errorf(string, ...interface{})
	Error(...interface{})
	Warningf(string, ...interface{})
	Warning(...interface{})
	Infof(string, ...interface{})
	Info(...interface{})
	Debugf(string, ...interface{})
	Debug(...interface{})
}

type defaultLog struct {
	*log.Logger
}

var logger = &defaultLog{Logger: log.New(os.Stderr, "rpc ", log.LstdFlags)}

func (l *defaultLog) Errorf(f string, v ...interface{}) {
	l.Printf("ERROR: "+f+"\n", v...)
}

func (l *defaultLog) Error(v ...interface{}) {
	l.Println(v...)
}

func (l *defaultLog) Warningf(f string, v ...interface{}) {
	l.Printf("WARNING: "+f+"\n", v...)
}

func (l *defaultLog) Warning(v ...interface{}) {
	l.Println(v...)
}

func (l *defaultLog) Infof(f string, v ...interface{}) {
	l.Printf("INFO: "+f+"\n", v...)
}

func (l *defaultLog) Info(v ...interface{}) {
	if IsDebug {
		l.Println(v...)
	}
}

func (l *defaultLog) Debugf(f string, v ...interface{}) {
	if IsDebug {
		l.Printf("DEBUG: "+f+"\n", v...)
	}
}

func (l *defaultLog) Debug(v ...interface{}) {
	if IsDebug {
		l.Println(v...)
	}
}
