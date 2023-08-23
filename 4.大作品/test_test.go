package main

import (
	"errors"
	"fmt"
	"testing"
)

var (
	test1Error = errors.New("test1")
	test2Error = errors.New("test2")
)



func TestError1(t *testing.T) {
	testMixError := fmt.Errorf("%w,123", test1Error)
	fmt.Println("testMixError", testMixError)
	fmt.Println("is test1:",errors.Is(testMixError, test1Error))
	fmt.Println("is test2:",errors.Is(testMixError, test2Error))
}


type Test3Error struct {
	s string
	status int
}

func (e *Test3Error) Error() string {
	return e.s
}

func newTest3Error(text string, status int) error {
	return &Test3Error{text,status}
}


func TestError2(t *testing.T) {
	//t1:= newTest3Error(fmt.Errorf("is test1:%w",test1Error).Error(),1)
	//t1:= newTest3Error(fmt.Errorf("is test3").Error(),1)
	t1:= newTest3Error("123",1)

	fmt.Println("t1",t1)
	var test3 *Test3Error
	fmt.Println("test3", test3)
	fmt.Println("is test3:", errors.As(t1, &test3))
	fmt.Printf("t1:%#v,test3:%#v",t1,test3)
}