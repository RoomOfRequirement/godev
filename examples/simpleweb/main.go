package main

import (
	"godev/examples/simpleweb/bb"
	"net/http"
)

// following geektutu's tutorials

func main() {
	b := bb.Default()

	// curl "http://localhost:6666/"
	b.GET("/", func(ctx *bb.Context) {
		ctx.HTML(http.StatusOK, "<h1>BB</h1>")
	})

	// curl "http://localhost:6666/v1/hello?name=test"
	g1 := b.AddGroup("/v1")
	g1.GET("/hello", func(ctx *bb.Context) {
		ctx.String(http.StatusOK, "hello %s", ctx.Query("name"))
	})

	g2 := b.AddGroup("/v2")
	// curl "http://localhost:6666/v2/login" -X POST -d 'username=test&password=000000'
	g2.POST("/login", func(ctx *bb.Context) {
		ctx.JSON(http.StatusOK, bb.OBJ{
			"username": ctx.PostForm("username"),
			"password": ctx.PostForm("password"),
		})
	})

	// curl "http://localhost:6666/panic/"
	g3 := b.AddGroup("/panic")
	g3.GET("/", func(ctx *bb.Context) {
		str := ""
		ctx.String(http.StatusOK, "%s", str[1])
	})

	_ = b.Run(":6666")
}
