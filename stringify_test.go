// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"encoding/json"
	"testing"
)

func TestStringify(t *testing.T) {
	doc := map[string]string{"color": "Red", "style": "Truck"}
	str := Stringify(doc)
	var v map[string]interface{}
	json.Unmarshal([]byte(str), &v)
	if v["color"] != doc["color"] || v["style"] != doc["style"] {
		t.Fatal("Expected", `{"color":"Red","style":"Truck"}`, "but got", str)
	}
}
