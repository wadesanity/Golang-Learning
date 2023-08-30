package main

import (
	"testing"
)

type NumAll interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | uint | uint8 | uint16 | uint32 | uint64
}

var aList [1]any

func Test2(t *testing.T) {
	aList[0] = any(1)
	// var a uint
	// a, ok := aList[0].(uint)
	// t.Log(a, ok)
	switch s := aList[0].(type) {
	case uint:
		t.Logf("uint:%T,%v", s, s)
	default:
		t.Logf("default:%T,%v,aList[0]:%T", s, s, aList[0])
	}
}
