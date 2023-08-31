package main

import "testing"

var a = 1

func Test1(t *testing.T) {
	switch a {
	case 1:
		t.Log("is 1")
	default:
		t.Log("no 1")
	}
}
