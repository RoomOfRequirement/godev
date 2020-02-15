package bb

import (
	"log"
	"time"
)

// Logger middleware
func Logger() HandleFunc {
	return func(ctx *Context) {
		start := time.Now()
		// process request
		ctx.Next()
		// calculate time consumption
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(start))
	}
}
