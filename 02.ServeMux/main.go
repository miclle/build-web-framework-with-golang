package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {

	// 这里生成一个 ServeMux 实例
	handler := http.NewServeMux()
	// handler := http.DefaultServeMux

	// 注册路由 1: /hello/
	handler.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.Replace(r.URL.Path, "/hello/", "", 1)

		io.WriteString(w, fmt.Sprintf("Hello %s\n", name))
	})

	// 注册路由 2: /hello
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "Hello, world!\n")
	})

	// 注册路由 3: /
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)

		io.WriteString(w, fmt.Sprintf("Oops Not found\nURL: %s\n", r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", handler))
}
