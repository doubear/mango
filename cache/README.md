# cache provider for mango

## Example

```go
m.SetCacher(cache.Memory(15*time.Minute))

func(ctx *mango.Context) (int, interface{}) {
    ctx.C.Set("cached-value", "cached value")
}
```
