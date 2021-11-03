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
size: 16:9
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

-----------------------------------------------------------------------

# Simple HTTP Server

使用 [`net/http#ListenAndServe`](https://golang.org/pkg/net/http/#ListenAndServe) 包实现一个最简单、最基础的 HTTP 服务。

```go
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Hello world, the web server
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

	http.HandleFunc("/hello", helloHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

```
go run main.go
```

-----------------------------------------------------------------------

#### 测试访问 /hello 路由：

```
curl -i 127.0.0.1:8080/hello
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
curl -i 127.0.0.1:8080/test
```

```
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 11 Dec 2020 12:09:55 GMT
Content-Length: 19

404 page not found
```

-----------------------------------------------------------------------

#### handler

handler 函数有两个参数，http.ResponseWriter 和 http.Request。 response writer 被用于写入 HTTP 响应数据，这里我们简单的返回 "Hello, world!\n"。

```go
helloHandler := func(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}
```

-----------------------------------------------------------------------

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

-----------------------------------------------------------------------

#### 两个结论：

1. `http.HandleFunc` 会将指定 `pattern` (模式、路由) 的 `handler` 注册在 `DefaultServeMux` 上面
2. `http.ListenAndServe` 如果 `handler` 为 `nil` ，在这种情况下使用 `DefaultServeMux` 。

-----------------------------------------------------------------------

#### 那么问题来了 `DefaultServeMux` 是啥？

```go
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}

type muxEntry struct {
	h       Handler
	pattern string
}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux { return new(ServeMux) }

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = &defaultServeMux

var defaultServeMux ServeMux

func (mux *ServeMux) Handle(pattern string, handler Handler)
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string)
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request)
```

-----------------------------------------------------------------------

# ServeMux

https://golang.org/pkg/net/http/#ServeMux

>ServeMux is an HTTP request multiplexer. It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.

**ServeMux 是一个 HTTP 请求多路复用器。它根据已注册模式列表匹配每个传入请求的 URL，并调用与 URL 最匹配的模式的处理程序。**

>Patterns name fixed, rooted paths, like "/favicon.ico", or rooted subtrees, like "/images/" (note the trailing slash). Longer patterns take precedence over shorter ones, so that if there are handlers registered for both "/images/" and "/images/thumbnails/", the latter handler will be called for paths beginning "/images/thumbnails/" and the former will receive requests for any other paths in the "/images/" subtree.

**匹配模式固定，较长的模式优先于较短的模式，"/" 匹配子树中任何其他路径的请求**

-----------------------------------------------------------------------

我们这次不使用默认的 `ServeMux` 来完成路由功能：

```go
func main() {
	// 这里生成一个 ServeMux 实例
	handler := http.NewServeMux()

	// 注册路由 /hello/
	handler.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.Replace(r.URL.Path, "/hello/", "", 1)
		io.WriteString(w, fmt.Sprintf("Hello %s\n", name))
	})

	// 注册路由 /hello
	handler.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	})

	// 注册路由 /
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, fmt.Sprintf("Oops Not found\nURL: %s\n", r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", handler))
}

```

-----------------------------------------------------------------------

#### 测试访问 /hello 路由：
```
curl -i 127.0.0.1:8080/hello
```

```
HTTP/1.1 200 OK
Date: Sun, 13 Dec 2020 08:46:13 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

Hello, world!
```

#### 测试访问 /hello/ 路由：

```
curl -i 127.0.0.1:8080/hello/
```
```
HTTP/1.1 200 OK
Date: Sun, 13 Dec 2020 08:46:16 GMT
Content-Length: 7
Content-Type: text/plain; charset=utf-8

Hello
```

-----------------------------------------------------------------------

#### 测试访问 /hello/foo 路由：

```
curl -i 127.0.0.1:8080/hello/foo
```
```
HTTP/1.1 200 OK
Date: Sun, 13 Dec 2020 08:48:17 GMT
Content-Length: 10
Content-Type: text/plain; charset=utf-8

Hello foo
```

#### 测试访问 /hello/foo/boo 路由：

```
curl -i 127.0.0.1:8080/hello/foo/boo
```
```
HTTP/1.1 200 OK
Date: Sun, 13 Dec 2020 08:48:30 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

Hello foo/boo
```

-----------------------------------------------------------------------

#### 测试访问不存在的 /test 路由：

```
curl -i 127.0.0.1:8080/test
```

```
HTTP/1.1 404 Not Found
Content-Type: text/plain
Date: Sun, 13 Dec 2020 09:11:45 GMT
Content-Length: 26

Oops Not found
URL: /test
```

#### 测试访问不存在的 /hel/foo 路由：

```
curl -i 127.0.0.1:8080/hel/foo
```

```
HTTP/1.1 404 Not Found
Content-Type: text/plain
Date: Sun, 13 Dec 2020 09:12:18 GMT
Content-Length: 29

Oops Not found
URL: /hel/foo
```

-----------------------------------------------------------------------

#### 这里发生两处变化：

1. 所有 `/hello/` 的子路径都被路由 1 接管，`/hello/` 后的子路径被赋值给 `name`。
2. 注册了 `/` 的路由，所以所有没有匹配到前两个路由的 URL 都会被路由 3 接管

<hr />

#### 默认的 DefaultServeMux 和自己定义的 ServeMux 对象有什么区别呢？

没有太大区别，完全可以把上面代码中的`handler := http.NewServeMux()` 这一行改为 `handler := http.DefaultServeMux`。
其实 `http.DefaultServeMux` 本身就是一个 `ServeMux` 类型的变量，只是为了方便，为 http 包添加必要的 API 提供了便利罢了。类似 log 包下的 `std`

```go
var std = New(os.Stderr, "", LstdFlags)

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	std.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}
```

-----------------------------------------------------------------------
## ServeMux 如何注册 handler？ `HandleFunc` 与 `Handle`

```go
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}
```

```go
func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	...

	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e // 放在 map 里
	if pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, e) // 排序后放在 slice 里
	}

	...
}

```

-----------------------------------------------------------------------

## ServeMux 如何匹配路由并分配处理器？

再回顾一下 `http.ListenAndServe` 的第二个参数：

```go
func ListenAndServe(addr string, handler Handler) error
```

Go 支持外部实现路由器，`ListenAndServe` 的第二个参数就是配置外部路由器，它是一个 `Handler` 接口。即外部路由器实现 `Hanlder` 接口。

https://golang.org/pkg/net/http/#Handler

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

ServeMux 实现了 ServeHTTP 方法

-----------------------------------------------------------------------

### ServeMux ServeHTTP

```go
// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
	h, _ := mux.Handler(r) // 找到对应的 Handler
	h.ServeHTTP(w, r)      // 响应请求
}
```

-----------------------------------------------------------------------

### ServeMux Handler

```go
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {

	...

	return mux.handler(host, r.URL.Path)
}

...

// handler is the main implementation of Handler.
// The path is known to be in canonical form, except for CONNECT methods.
func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path)
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}

```

-----------------------------------------------------------------------

```go
// Find a handler on a handler map given a path string.
// Most-specific (longest) pattern wins.
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	// Check for exact match first.
	v, ok := mux.m[path]
	if ok {
		return v.h, v.pattern
	}

	// Check for longest valid match.  mux.es contains all patterns
	// that end in / sorted from longest to shortest.
	for _, e := range mux.es {
		if strings.HasPrefix(path, e.pattern) {
			return e.h, e.pattern
		}
	}
	return nil, ""
}
```

-----------------------------------------------------------------------

## ServeMux 路由器设计思路

#### 注册路由：
1. 使用 `ServeMux.HandlerFunc` 注册 `func(ResponseWriter, *Request)` 签名的函数作为处理器：
1.1 在内部转换为 `http.HandlerFunc` 对象，`http.HandlerFunc` 类型实现了 `http.Handler` 接口
1.2 之后再调用 `ServeMux.Handle` 方法注册路由
2. 使用 `ServeMuxHandle` 注册 `http.Handler` 对象作为处理器
2.1 将 handler 保存在 ServeMux 内置的 muxEntry map 和 slice 中

#### 匹配并处理路由：
1. 通过 `http.ListenAndServe(addr, mux)` ServeMux.ServeHTTP 接收请求
2. 使用 ServeMux.Handler 匹配合适路由，并返回 handler
2.1 ServeMux.Handler -> ServeMux.handler(host, r.URL.Path)
2.2 ServeMux.handler(host, r.URL.Path) -> ServeMux.match(host + path | path) 匹配路由
3. 调用 handler.ServeHTTP(w, r) 处理请求

-----------------------------------------------------------------------

## 扩展阅读：一个 HTTP 连接处理流程

![An HTTP connection processing flow](./assets/http-connection-processing-flow.png)

https://astaxie.gitbooks.io/build-web-application-with-golang/content/zh/03.3.html

