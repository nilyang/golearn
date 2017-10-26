// github.com/gorilla/websocket 示例 练习&理解

package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8888", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "No found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	http.ServeFile(w, r, "home.html")
}

func main() {
	runHttp()
}


func runHttp() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 该闭包每次请求时，运行在新启用的 goroutine 里
		serveWebsocket(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

