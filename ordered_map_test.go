// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"encoding/json"
	"testing"
)

func TestNewOrderedMap(t *testing.T) {
	str := `{"color":"Red","style":"Truck","brand":"BMW","child":{"brand":"BMW","color":"Red","style":"Truck"}}`
	om := NewOrderedMap(str)
	data, _ := json.Marshal(om)
	if string(data) != str {
		t.Fatal("Expected", str, "but got", string(data))
	}
}
