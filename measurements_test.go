// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"testing"
)

func TestGetStorageSize(t *testing.T) {
	size := 123456789
	if GetStorageSize(size) != "117.7 MB" {
		t.Fatal()
	}
}

func TestGetDurationFromSeconds(t *testing.T) {
	seconds := float64(68)
	if GetDurationFromSeconds(seconds) != "1.1 minutes" {
		t.Fatal()
	}
}
