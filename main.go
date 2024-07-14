package main

import (
	"log"
	"net/http"
	"pong/auth"
	"pong/boundary"
	"pong/controller"
	"pong/repository"
	"runtime/debug"
	"time"

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

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("static/")))
	mux.Handle("POST /user", userHandler)
	mux.Handle("/match", auth.BasicAuth(matchHandler, &userStore))

	handler := Logging(mux)
	handler = PanicRecovery(handler)
	log.Fatalln(http.ListenAndServe(":"+Port, handler))
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}
