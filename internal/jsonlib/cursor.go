// Copyright 2020 Navibyte (https://navibyte.com). All rights reserved.
// Use of this source code is governed by a MIT-style license, see the LICENSE.

package jsonlib

import (
	"github.com/tidwall/gjson"
)

// Cursor is a wrapper for gjson.Result providing easier interface for getting
// values from a JSON string in some use cases (ie. copying values on different
// types from JSON to structs). The wrapper may have some performance penalty
// compared to using directly gjson.Result but on the other hand it allows
// (domain specific) parsing code to be a bit less verbose.
type Cursor struct {
	result gjson.Result
}

// NewCursor creates a new cursor wrapping gjson.Result.
func NewCursor(result gjson.Result) Cursor {
	return Cursor{result: result}
}

// Result returns the gjson.Result the cursor is pointing at.
func (c Cursor) Result() gjson.Result {
	return c.result
}

// Get returns a new cursor for a given path.
func (c Cursor) Get(path string) Cursor {
	return Cursor{result: c.result.Get(path)}
}

// Go moves the cursor to a new path location if such location exists (and
// returning true). If a new path location does not exits, returns false.
func (c Cursor) Go(path string) bool {
	newResult := c.result.Get(path)
	if !newResult.Exists() {
		return false
	}
	c.result = newResult
	return true
}

// ForEachArray iterates through array on a path location.
func (c Cursor) ForEachArray(path string, iterator func(value Cursor) bool) {
	if array := c.result.Get(path); array.IsArray() {
		array.ForEach(func(_ gjson.Result, v gjson.Result) bool {
			iterator(NewCursor(v))
			return true
		})
	}
}

// Exists returns true if the cursor points to an existing JSON value.
func (c Cursor) Exists() bool {
	return c.result.Exists()
}

// IsObject returns true if the cursor points to a JSON object.
func (c Cursor) IsObject() bool {
	return c.result.IsObject()
}

// IsArray returns true if the cursor points to a JSON array.
func (c Cursor) IsArray() bool {
	return c.result.IsArray()
}

// String returns an element as a string.
func (c Cursor) String(path string) string {
	return c.result.Get(path).String()
}

// Bool returns an element as a bool.
func (c Cursor) Bool(path string) bool {
	return c.result.Get(path).Bool()
}

// Int32 returns an element as an int32.
func (c Cursor) Int32(path string) int32 {
	return int32(c.result.Get(path).Int())
}

// Int64 returns an element as an int64.
func (c Cursor) Int64(path string) int64 {
	return c.result.Get(path).Int()
}

// Float32 returns an element as an float32.
func (c Cursor) Float32(path string) float32 {
	return float32(c.result.Get(path).Float())
}

// Float64 returns an element as an float64.
func (c Cursor) Float64(path string) float64 {
	return c.result.Get(path).Float()
}
