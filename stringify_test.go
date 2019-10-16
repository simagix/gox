// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"encoding/json"
	"testing"
)

func TestStringify(t *testing.T) {
	doc := map[string]interface{}{"color": "Red", "style": "Truck", "amount": 637.32}
	str := Stringify(doc)
	var v map[string]interface{}
	json.Unmarshal([]byte(str), &v)
	if v["color"] != doc["color"] || v["style"] != doc["style"] || v["amount"] != doc["amount"] {
		t.Fatal("Expected", `{"color":"Red","style":"Truck","amount":637.32}`, "but got", str)
	}
}
