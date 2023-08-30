package main

import (
	"testing"
)

type User struct {
	Id   uint
	Name string
	Age  uint
}

type APIUser struct {
	Id   uint
	Name string
}

type Video struct {
	Id    uint
	Title string
	path  string
}

type ApiVideo struct {
	Id    uint
	Title string
}

type MO interface {
	User | Video
}

type TO interface {
	APIUser | ApiVideo
}

func Test3(t *testing.T) {
	r := f1(&User{}, &APIUser{})
	r.Id = 123
	r.Name = "haha"
	t.Errorf("r:%#v", r)
}

func f1[M MO, T TO](m *M, t *T) *T {
	// fmt.Printf("m:%T, m:%#v", m, m)
	// fmt.Printf("t:%T, t:%#v", t, t)
	return t
}
