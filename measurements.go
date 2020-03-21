// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"fmt"
	"strconv"
)

// GetStorageSize returns storage size in [TGMK] B
func GetStorageSize(num interface{}) string {
	f := fmt.Sprintf("%v", num)
	x, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return f
	}

	if x >= (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fTB", x/(1024*1024*1024*1024))
	} else if x >= (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fGB", x/(1024*1024*1024))
	} else if x >= (1024 * 1024) {
		return fmt.Sprintf("%.1fMB", x/(1024*1024))
	} else if x >= 1024 {
		return fmt.Sprintf("%.1fKB", x/1024)
	}
	return fmt.Sprintf("%v", x)
}

// GetDurationFromSeconds converts seconds to time string, e.g. 1.5m
func GetDurationFromSeconds(seconds float64) string {
	timestr := fmt.Sprintf("%3.0f", seconds)
	if seconds >= (24 * 60 * 60) {
		seconds /= (24 * 60 * 60)
		timestr = fmt.Sprintf("%4.1f days", seconds)
	} else if seconds >= (60 * 60) {
		seconds /= (60 * 60)
		timestr = fmt.Sprintf("%3.1f hours", seconds)
	} else if seconds >= 60 {
		seconds /= 60
		timestr = fmt.Sprintf("%3.1f minutes", seconds)
	} else if seconds >= 1 {
		timestr = fmt.Sprintf("%3.0f secnods", seconds)
	}
	return timestr
}
