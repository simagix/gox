// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {
	docMap := make(map[string]interface{})
	docMap["color"] = "Red"
	docMap["array"] = []string{"hello", "world"}
	docMap["colors"] = []string{"Red", "Yellow", "Brown", "Black", "White"}
	length := len(docMap["colors"].([]string))
	walker := NewMapWalker(cb)
	doc := walker.Walk(docMap)
	if doc.(map[string]interface{})["color"] != "Green" {
		t.Fatal("Expected", "Green", "but got", doc.(map[string]interface{})["color"])
	}
	if maxArrayLength != length {
		t.Fatal("Expected", length, "but got", maxArrayLength)
	}
	t.Log(Stringify(doc, "", "  "))

	map2 := map[string]interface{}{}
	map2["doc"] = doc
	doc = walker.Walk(map2)
	if walker.GetNestedLevel() != 2 {
		t.Fatal(walker.GetNestedLevel())
	}
	t.Log(Stringify(doc, "", "  "))

	map3 := map[string]interface{}{}
	map3["three"] = map2
	map3["four"] = map2
	doc = walker.Walk(map3)
	t.Log(Stringify(doc, "", "  "))
	if walker.GetNestedLevel() != 3 {
		t.Fatal(walker.GetNestedLevel())
	}
}

var maxArrayLength = 0

func cb(value interface{}) interface{} {
	if v, ok := value.(string); ok {
		if v == "Red" {
			return "Green"
		}
	} else if reflect.TypeOf(value).Kind() == reflect.Array ||
		reflect.TypeOf(value).Kind() == reflect.Slice {
		length := len(value.([]interface{}))
		if length > maxArrayLength {
			maxArrayLength = length
		}
	}
	return value
}
