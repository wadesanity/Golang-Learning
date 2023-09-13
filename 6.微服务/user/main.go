package main

import (
	"user/grpc/serve"
	"user/repository/db"
)

func main() {
	db.Init()
	serve.Start()
}
