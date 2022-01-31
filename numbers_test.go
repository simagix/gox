// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import "testing"

func TestToInt64(t *testing.T) {
	var err error
	var x int64
	numbers := map[string]interface{}{"int": 123, "float64": 123.45}
	if x, err = ToInt64(numbers["int"]); err != nil {
		t.Fatal(err)
	}
	assertEqual(t, int64(123), x)
	if x, err = ToInt64(numbers["float64"]); err != nil {
		t.Fatal(err)
	}
	assertEqual(t, int64(123), x)
	if _, err = ToInt64(numbers["nil"]); err == nil {
		t.Fatal("expects error")
	}
	t.Log(err)
}

func TestToFloat64(t *testing.T) {
	var err error
	var x float64
	numbers := map[string]interface{}{"int": 123, "float64": 123.45}
	if x, err = ToFloat64(numbers["int"]); err != nil {
		t.Fatal(err)
	}
	assertEqual(t, float64(123), x)
	if x, err = ToFloat64(numbers["float64"]); err != nil {
		t.Fatal(err)
	}
	assertEqual(t, 123.45, x)
	if _, err = ToFloat64(numbers["nil"]); err == nil {
		t.Fatal("expects error")
	}
	t.Log(err)
}
