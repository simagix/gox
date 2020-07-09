// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import "testing"

func TestToInt64(t *testing.T) {
	var err error
	var x int64
	numbers := map[string]interface{}{"int": 123, "float64": 123.45}
	x, err = ToInt64(numbers["int"])
	if x != 123 {
		t.Fatal("expects", 123, "but had", x)
	}
	x, err = ToInt64(numbers["float64"])
	if x != 123 {
		t.Fatal("expects", 123, "but had", x)
	}
	if x, err = ToInt64(numbers["nil"]); err == nil {
		t.Fatal("expects error")
	}
	t.Log(err)
}

func TestToFloat64(t *testing.T) {
	var err error
	var x float64
	numbers := map[string]interface{}{"int": 123, "float64": 123.45}
	x, err = ToFloat64(numbers["int"])
	if x != 123 {
		t.Fatal("expects", 123, "but had", x)
	}
	x, err = ToFloat64(numbers["float64"])
	if x != 123.45 {
		t.Fatal("expects", 123, "but had", x)
	}
	if x, err = ToFloat64(numbers["nil"]); err == nil {
		t.Fatal("expects error")
	}
	t.Log(err)
}
