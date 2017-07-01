# JWT middleware of Mango

## Usage

```go
//create jwt instance.
jwtt = jwt.New(jwt.HS256, "secret token")

//set errors handler
jwtt.Error(func(e error, ctx *mango.Context) {
    ctx.W.Clear()

    switch e {
    case jwt.ErrTokenExpired:
        ctx.W.WriteJSON(map[string]interface{}{
            "status":  "401",
            "message": "Token is expired.",
        })
    case jwt.ErrTokenInvalid:
        ctx.W.WriteJSON(map[string]interface{}{
            "status":  "401",
            "message": "Token is invalid.",
        })
    case jwt.ErrTokenLost:
        ctx.W.WriteJSON(map[string]interface{}{
            "status":  "401",
            "message": "Token is lost.",
        })
    default:
        ctx.W.SetStatus(500)
    }
})

//use middleware generator to generates jwt auto validate middleware.
m.Get("/index", index, jwtt.Auth("porta"))

m.Post("/auth", func(ctx *mango.Context) (int, interface{}) {

    //create claims instance with user credential.
	c := jwt.Claims{
		Sub: "username", //user id or other credential.
		Aud: "portal",
	}

    //call Sign to create token.
	t := jwtt.Sign(c)

	return 200, map[string]interface{}{
		"token": t,
	}
})
```
