// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"testing"
)

func TestWalk(t *testing.T) {
	docMap := make(map[string]interface{})
	docMap["color"] = "Red"
	walker := NewMapWalker(cb)
	doc := walker.Walk(docMap)
	if doc["color"] != "Green" {
		t.Fatal("Expected", "Green", "but got", doc["color"])
	}
}

func cb(value interface{}) interface{} {
	if v, ok := value.(string); ok {
		if v == "Red" {
			return "Green"
		}
	}
	return value
}
