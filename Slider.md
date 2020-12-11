---
theme: gaia
style: |
  section {
    background-color: #ccc;
    font-size: 28px;
    padding: 20px 24px 24px 24px;
  }
  pre {
    margin: 0.5em 0 0 0;
  }
size: 4K
_class: lead
paginate: true
backgroundColor: #fff
backgroundImage: url('./assets/background.jpg')
marp: true
---

<!-- ![bg left:40% 80%](https://marp.app/assets/marp.svg) -->

# **Build web framework with golang**

Miclle Zheng

@miclle

---

# Simple HTTP Server

使用 `net/http` 包实现一个最简单、最基础的 HTTP 服务。

```go
package main

import (
	"io"
	"net/http"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func main() {

	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":3000", nil)
}
```

```
go run main.go
```

---

#### 测试访问 /hello 路由：

```
curl -i 127.0.0.1:3000/hello
```

```
HTTP/1.1 200 OK
Date: Fri, 11 Dec 2020 12:09:21 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

Hello, world!
```

#### 测试访问不存在的 /test 路由：

```
curl -i 127.0.0.1:3000/test
```

```
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 11 Dec 2020 12:09:55 GMT
Content-Length: 19

404 page not found
```

---

#### handler

`handlers` 是 `net/http` 服务器里面的一个基本概念。
`handler` 对象实现了 `http.Handler` 接口。
编写 `handler` 的常见方法是，在具有适当签名的函数上使用 `http.HandlerFunc` 适配器。

handler 函数有两个参数，http.ResponseWriter 和 http.Request。 response writer 被用于写入 HTTP 响应数据，这里我们简单的返回 "Hello, world!\n"。

```go
func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}
```

---

#### 源码分析

http.HandleFunc
```go
// HandleFunc registers the handler function for the given pattern
// in the DefaultServeMux.
// The documentation for ServeMux explains how patterns are matched.
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

http.ListenAndServe
```go
// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```

---

#### 两个结论：

1. `http.HandleFunc` 会将指定 `pattern` (模式、路由) 的 `handler` 注册在 `DefaultServeMux` 上面
2. `http.ListenAndServe` 如果 `handler` 为 `nil` ，在这种情况下使用 `DefaultServeMux` 。


#### 那么问题来了 `DefaultServeMux` 是啥？
不管是啥，肯定是实现了 http.Handler 接口

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```