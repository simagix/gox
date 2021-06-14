// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"testing"
)

func TestInfo(t *testing.T) {
	logger := GetLogger("TestInfo")
	logger.Info("first ", "second ", 3, 4)
}

func TestInfof(t *testing.T) {
	logger := GetLogger("TestInfof")
	logger.Infof(`%v => %v`, "key", 100)
}

func TestErrorf(t *testing.T) {
	logger := GetLogger("TestErrorf")
	logger.Errorf(`%v => %v`, "key", 100)
}

func TestWarnf(t *testing.T) {
	logger := GetLogger("TestWarnf")
	logger.Warnf(`%v => %v`, "key", 100)
}

func TestDebug(t *testing.T) {
	logger := GetLogger("TestDebug")
	logger.Debug(`no message`)
	logger.SetLoggerLevel(Debug)
	logger.Debug(`debug message`)
	logger.SetLoggerLevel(Info)
	logger.Debug(`no message`)
	logger.SetLoggerLevel(Trace)
	logger.Debug(`trace message`)
}
