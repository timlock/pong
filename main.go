package main

import (
	"log"
	"net/http"
	"pong/boundary"
	"pong/controller"
	"pong/repository"

	"github.com/gorilla/websocket"
)

const (
	Port  = "8080"
	Ticks = 33
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	userStore := repository.NewUserStore()
	simulator := controller.NewSimulator(Ticks)
	simulator.Run()
	userHandler := boundary.NewUserHandler(&userStore)
	matchHandler := boundary.NewMatchHanlder(&userStore, &simulator, upgrader)

	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.Handle("/user", userHandler)
	http.Handle("/match", matchHandler)
	log.Fatalln(http.ListenAndServe(":"+Port, nil))
}
