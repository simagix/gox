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
	num := len(logger.Logs)

	logger.Debug(`no message`)
	t.Log(num, len(logger.Logs))
	assertEqual(t, num, len(logger.Logs))

	logger.SetLoggerLevel(Debug)
	logger.Debug(`debug message`)
	t.Log(num, len(logger.Logs))
	assertEqual(t, num, len(logger.Logs))

	logger.SetLoggerLevel(Info)
	logger.Debug(`no message`)
	t.Log(num, len(logger.Logs))
	assertEqual(t, num, len(logger.Logs))

	logger.SetLoggerLevel(Trace)
	logger.Debug(`trace message`)
	t.Log(num, len(logger.Logs))
	assertEqual(t, num, len(logger.Logs))
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger()
	assertEqual(t, "gox", logger.AppName)
	assertEqual(t, true, logger.inMem)
}

func TestGetLoggerOneParam(t *testing.T) {
	appName := "myapp"
	logger := GetLogger(appName)
	assertEqual(t, appName, logger.AppName)
	assertEqual(t, true, logger.inMem)
}

func TestGetLoggerTwoParams(t *testing.T) {
	appName := "myapp"
	logger := GetLogger(appName, false)
	assertEqual(t, appName, logger.AppName)
	assertEqual(t, false, logger.inMem)

	logger.Info(appName)
	logger.Debug(appName)
	assertEqual(t, 1, len(logger.Logs))
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
