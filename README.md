# Mango - keep simple, keep easy to use.
[![GoDoc](https://godoc.org/github.com/go-mango/mango?status.svg)](https://godoc.org/github.com/go-mango/mango)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mango/mango)](https://goreportcard.com/report/github.com/go-mango/mango)

## Quick Start

```go
func index(ctx *mango.Context) (int, interface{}) {
	return 200, map[string]interface{}{
		"message": "hello mango",
	}
}

func main() {
	m := mango.Default()
	m.Get("/", index)
	m.Start(":8080")
}
```

## Route Definition

### Basic Guide

```go
m.Get("/", index)
m.Post("/", index)
m.Put("/", index)
m.Delete("/", index)
m.Any("/any", index) //GET,POST,PUT,DELETE
```

### Routes Group

```go
m.Group("/api", func(api *mango.GroupRouter) {
	api.Get("/", getApi)
	api.Post("/", postApi)
	api.Put("/", putApi)
	api.Delete("/", deleteApi)
	api.Any("/", anyApi) //GET,POST,PUT,DELETE
})
```

### Nested Routes Group

```go
m.Group("/api", func(api *mango.GroupRouter) {
	api.Group("/v1", func(v1 *mango.GroupRouter) {
		v1.Get("/", getApiV1)
		v1.Post("/", postApiV1)
		v1.Put("/", putApiV1)
		v1.Delete("/", deleteApiV1)
		v1.Any("/", anyApiV1) //GET,POST,PUT,DELETE
	})
})
```

### Route Middleware

```go
func myMiddleware() mango.MiddlerFunc {
	return func(ctx *mango.Context) {

		//do something before route handler executed.

		ctx.Next()

		//do something before response sent to client.
	}
}

m.Use(myMiddleware())

//or

m.Get("/", index, myMiddleware())

//or

m.Group("/api", func(api *mango.GroupRouter){

}, myMiddleware())
```

## Built-in Middlewares

1. Record
2. Recovery
3. Cors
4. Static
5. Redirect
6. Compress
7. ...

## Serve Mode

### HTTP Mode

```go
m.Start(":8080")
```

### HTTPS Mode

```go
m.StartTLS(":8080", "cert file content", "key file content")
```

### HTTPS with Let's encrypt Mode

```go
m.StartAutoTLS(":8080", "example.org")
```

### Handle HTTP and HTTPS

```go
go m.Start(":http")
m.StartAutoTLS(":https", "example.org")
```

## Benchmark
??? what's that???

## License
Based on MIT license.