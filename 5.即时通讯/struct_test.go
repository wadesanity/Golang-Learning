package main

import (
	"fmt"
	"testing"
)

type a struct {
	f1 string
}

func Test1(t *testing.T) {
	var a1 a
	a1.f1 = "1"
	fmt.Printf("v:%#v, addr:%p", a1, &a1)
	a1.f1 = "2"
	fmt.Printf("v:%#v, addr:%p", a1, &a1)
}
