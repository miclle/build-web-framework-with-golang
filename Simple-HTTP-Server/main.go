package main

import (
	"io"
	"net/http"
)

// `handlers` 是 `net/http` 服务器里面的一个基本概念。
// `handler` 对象实现了 `http.Handler` 接口。
// 编写 `handler` 的常见方法是，在具有适当签名的函数上使用 `http.HandlerFunc` 适配器。
func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func main() {

	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":3000", nil)
}
