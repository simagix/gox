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
	if doc.(map[string]interface{})["color"] != "Green" {
		t.Fatal("Expected", "Green", "but got", doc.(map[string]interface{})["color"])
	}

	walker = NewMapWalker(nil)
	map2 := map[string]interface{}{}
	map2["doc"] = doc
	doc = walker.Walk(map2)
	t.Log(doc)
	if walker.GetNestedLevel() != 2 {
		t.Fatal(walker.GetNestedLevel())
	}

	walker = NewMapWalker(nil)
	map3 := map[string]interface{}{}
	map3["three"] = map2
	map3["four"] = map2
	doc = walker.Walk(map3)
	t.Log(doc)
	if walker.GetNestedLevel() != 3 {
		t.Fatal(walker.GetNestedLevel())
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
