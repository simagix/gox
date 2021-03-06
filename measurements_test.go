// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"testing"
)

func TestGetStorageSize(t *testing.T) {
	size := float64(123456789)
	if s := GetStorageSize(size); s != "117.7MB" {
		t.Fatal(s)
	}

	size = 3.14159262
	if s := GetStorageSize(size); s != "3" {
		t.Fatal(s)
	}
}

func TestGetDurationFromSeconds(t *testing.T) {
	seconds := float64(68)
	if GetDurationFromSeconds(seconds) != "1.1 minutes" {
		t.Fatal()
	}
}
