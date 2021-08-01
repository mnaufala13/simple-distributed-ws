package main

import (
	"flag"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	rdsAddr := os.Getenv("REDIS_ADDR")
	if rdsAddr == "" {
		rdsAddr = "localhost:6379"
	}
	flag.Parse()
	rds := redis.NewClient(&redis.Options{
		Addr:     rdsAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	hub := newHub(rds)
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
