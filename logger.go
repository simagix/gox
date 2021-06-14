// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// specify logger level
const (
	Trace int = iota
	Debug
	Info
	Warn
	Error
)

// Logger stores logger info
type Logger struct {
	logs  []string
	level int
}

var instance *Logger
var once sync.Once

// GetLogger returns Logger instance
func GetLogger(appName string) *Logger {
	once.Do(func() {
		instance = &Logger{level: Info}
		instance.Infof(`%v begins at %v`, appName, time.Now().Format(time.RFC3339))
	})
	return instance
}

// SetLoggerLevel sets logger level
func (p *Logger) SetLoggerLevel(level int) {
	p.level = level
}

// Error adds and prints an error message
func (p *Logger) Error(v ...interface{}) {
	p.print("E", fmt.Sprint(v...), Error)
}

// Errorf adds and prints a message
func (p *Logger) Errorf(format string, v ...interface{}) {
	p.print("E", fmt.Sprintf(format, v...), Error)
}

// Warn adds and prints a warning message
func (p *Logger) Warn(v ...interface{}) {
	p.print("W", fmt.Sprint(v...), Warn)
}

// Warnf adds and prints a message
func (p *Logger) Warnf(format string, v ...interface{}) {
	p.print("W", fmt.Sprintf(format, v...), Warn)
}

// Info adds and prints a message
func (p *Logger) Info(v ...interface{}) {
	p.print("I", fmt.Sprint(v...), Info)
}

// Infof adds and prints a message
func (p *Logger) Infof(format string, v ...interface{}) {
	p.print("I", fmt.Sprintf(format, v...), Info)
}

// Debug adds and prints a message
func (p *Logger) Debug(v ...interface{}) {
	p.print("D", fmt.Sprint(v...), Debug)
}

// Debugf adds and prints a message
func (p *Logger) Debugf(format string, v ...interface{}) {
	p.print("D", fmt.Sprintf(format, v...), Debug)
}

// Trace adds and prints a message
func (p *Logger) Trace(v ...interface{}) {
	p.print("T", fmt.Sprint(v...), Trace)
}

// Tracef adds and prints a message
func (p *Logger) Tracef(format string, v ...interface{}) {
	p.print("T", fmt.Sprintf(format, v...), Trace)
}

// Log adds and prints a message
func (p *Logger) print(indicator string, message string, level int) {
	if level < p.level {
		return
	}
	str := fmt.Sprintf(`%v %v %v`, time.Now().Format(time.RFC3339), indicator, message)
	p.logs = append(p.logs, str)
	fmt.Println(str)
}

// Print prints all logs
func (p *Logger) Print() string {
	return strings.Join(p.logs, "\n")
}
