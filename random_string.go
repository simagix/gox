// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GetRandomDigitString returns a random digit string
func GetRandomDigitString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return string(bytes)
	}
	return string(bytes)
}

// GetRandomHexString returns a random hex string
func GetRandomHexString(n int) string {
	bytes := []byte(GetRandomDigitString(n / 2))
	return hex.EncodeToString(bytes)
}

// GetRandomUUIDString returns a random UUID string
func GetRandomUUIDString() string {
	return fmt.Sprintf(`UUID("%s-%s-%s-%s-%s")`, GetRandomHexString(8), GetRandomHexString(4), GetRandomHexString(4),
		GetRandomHexString(4), GetRandomHexString(12))
}
