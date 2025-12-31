// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"fmt"
	"strconv"
)

// ToInt converts any data type to int (ignores errors, returns 0 on failure)
func ToInt(num interface{}) int {
	x, _ := ToInt64(num)
	return int(x)
}

// ToInt64 converts any data type to int64
func ToInt64(num interface{}) (int64, error) {
	var err error
	var x float64
	if x, err = ToFloat64(num); err != nil {
		return 0, err
	}
	return int64(x), err
}

// ToFloat64 converts any data type to float64
func ToFloat64(num interface{}) (float64, error) {
	f := fmt.Sprintf("%v", num)
	return strconv.ParseFloat(f, 64)
}
