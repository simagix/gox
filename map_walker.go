// Copyright 2020 Kuei-chun Chen. All rights reserved.

package gox

import (
	"encoding/json"
	"reflect"
)

type callback func(interface{}) interface{}

// MapWalker is an empty JSON document
type MapWalker struct {
	arrayLength int
	cb          callback
	level       int
	nestedLevel int
}

// NewMapWalker returns a MapWalker
func NewMapWalker(cb callback) *MapWalker {
	return &MapWalker{cb: cb, nestedLevel: 1}
}

// SetCallBack defines callback function
func (walker *MapWalker) SetCallBack(cb callback) { walker.cb = cb }

// GetNestedLevel return the level of nested document
func (walker *MapWalker) GetNestedLevel() int { return walker.nestedLevel }

// Walk walks a map
func (walker *MapWalker) Walk(v interface{}) interface{} {
	walker.arrayLength = 0
	walker.level = 0
	walker.nestedLevel = 0
	return walker.traverse(v)
}

// traverse walks a map
func (walker *MapWalker) traverse(v interface{}) interface{} {
	if v == nil {
		return v
	}
	vt := reflect.TypeOf(v)
	if vt == nil {
		return v
	}
	switch vt.Kind() {
	case reflect.Map:
		vmap, ok := v.(map[string]interface{})
		if !ok {
			buf, err := json.Marshal(v)
			if err != nil {
				return v
			}
			json.Unmarshal(buf, &vmap)
		}
		walker.level++
		if walker.level > walker.nestedLevel {
			walker.nestedLevel = walker.level
		}
		for k, val := range vmap {
			vmap[k] = walker.traverse(val)
		}
		walker.level--
		return vmap
	case reflect.Array, reflect.Slice:
		arr, ok := v.([]interface{})
		if !ok {
			buf, err := json.Marshal(v)
			if err != nil {
				return v
			}
			json.Unmarshal(buf, &arr)
		}
		if len(arr) > walker.arrayLength {
			walker.arrayLength = len(arr)
		}
		for i, val := range arr {
			arr[i] = walker.traverse(val)
		}
		return arr
	}
	if walker.cb != nil {
		return walker.cb(v)
	}
	return v
}
