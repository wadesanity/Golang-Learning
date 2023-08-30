package main

import (
	"fmt"
	"testing"
)

type NumStr interface {
	Num | Str
}

type Num interface {
	~int | ~int32 | ~uint64
}

type Str interface {
	string
}

func add[T NumStr](a, b T) T {
	return a + b
}

func Test1(t *testing.T) {
	fmt.Println(add(3, 4))
	fmt.Println(add("dudu", "yiyi"))
}
