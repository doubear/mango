# Mango - keep simple, keep easy to use.

```go
package main

import "github.com/doubear/mango"

//Index index of system
func Index(ctx *mango.Context) {
	ctx.W.Write([]byte("hello world"))
}

func Show(ctx *mango.Context) {
	id := ctx.Param("id", "")
	ctx.W.WriteJSON(map[string]interface{}{
		"id": id,
	})
}

func ApiV1Show(ctx *mango.Context) {
	ctx.W.Write([]byte("API.v1"))
}

func ApiV2Show(ctx *mango.Context) {
	ctx.W.Write([]byte("hello API.v2"))
}

func main() {
	m := mango.New()
	m.Use(mango.Recovery())
	m.Use(mango.Record())
	m.Get("/", Index)
	m.Group("api", func(s1 *mango.GroupRouter) {
		s1.Group("v1", func(s2 *mango.GroupRouter) {
			s2.Get("/", ApiV1Show)
		})

		s1.Group("v2", func(s2 *mango.GroupRouter) {
			s2.Get("/", ApiV2Show)
		})
	})

	m.Get("/{id}", Show)
	m.Start("127.0.0.1:9000")
}
```