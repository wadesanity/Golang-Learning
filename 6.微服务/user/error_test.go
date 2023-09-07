package main

import (
	"fmt"
	"testing"
)

type Error string

func (e Error) Error() string { return string(e) }

func Test1(t *testing.T) {
	const err = Error("EOF")
	//const err2 = errors.errorString{"EOF"} // const initializer errorString literal is not a constant}
	//err = Error("not EOF") // error, cannot assign to err
	fmt.Println(err == Error("EOF")) // true
}
