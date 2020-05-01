// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import "testing"

func TestGetRandomDigitString(t *testing.T) {
	if n := len(GetRandomDigitString(16)); n != 16 {
		t.Fatal(n)
	}
}

func TestGetRandomHexString(t *testing.T) {
	if n := len(GetRandomHexString(16)); n != 16 {
		t.Fatal(n)
	}
}

func TestGetRandomUUIDString(t *testing.T) {
	if n := len(GetRandomUUIDString()); n != 44 {
		t.Fatal(n)
	}
}
