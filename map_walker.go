// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"reflect"
)

type callback func(interface{}) interface{}

// MapWalker is an empty JSON document
type MapWalker struct {
	level       int
	nestedLevel int
	cb          callback
}

// NewMapWalker returns a MapWalker
func NewMapWalker(cb callback) *MapWalker {
	return &MapWalker{cb: cb, nestedLevel: 1}
}

// GetNestedLevel return the level of nested document
func (walker *MapWalker) GetNestedLevel() int { return walker.nestedLevel }

// Walk walks a map
func (walker *MapWalker) Walk(v interface{}) interface{} {
	vt := reflect.TypeOf(v)
	switch vt.Kind() {
	case reflect.Map:
		if vmap, ok := v.(map[string]interface{}); ok {
			walker.level++
			if walker.level > walker.nestedLevel {
				walker.nestedLevel = walker.level
			}
			for k, val := range vmap {
				vmap[k] = walker.Walk(val)
			}
			walker.level--
			return vmap
		}
		return v
	case reflect.Array, reflect.Slice:
		if arr, ok := v.([]interface{}); ok {
			for i, val := range arr {
				arr[i] = walker.Walk(val)
			}
			return arr
		}
	default:
		if walker.cb != nil {
			return walker.cb(v)
		}
	}
	return v
}
