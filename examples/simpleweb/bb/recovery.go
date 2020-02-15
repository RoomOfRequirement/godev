package bb

import (
	"fmt"
	"godev/utils"
	"log"
	"net/http"
)

// Recovery middleware
func Recovery() HandleFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", utils.Trace(message))
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		// pass control to next handler
		ctx.Next()
	}
}
