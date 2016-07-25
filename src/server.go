package main

import (
	"runtime"
	"./socketserver"
	"./httpserver"
	"log"
)

const (
	SERVER_NETWORK = "tcp"
	SERVER_ADDRESS = "127.0.0.1:9988"
	HTTP_ADDRESS = "127.0.0.1:8000"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag := make(chan bool,1)
	go httpserver.StartHttpServer(HTTP_ADDRESS, flag)
	go socketserver.StartSocket(SERVER_NETWORK,SERVER_ADDRESS, flag)
	<-flag
	log.Fatal("socket server stop......................")
}

