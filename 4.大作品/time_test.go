package main

import (
	"fmt"
	"testing"
	"time"
)

func Test4(t *testing.T) {
	unixTime := time.Now()
	fmt.Println(unixTime.Unix())
}
