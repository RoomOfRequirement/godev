package bb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context ...
type Context struct {
	// parameters for serve http
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Method string
	Path   string
	// response info
	StatusCode int
	// middleware
	handlers []HandleFunc
	idx      int // handler execution cursor
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:     w,
		Req:        r,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: 0,
		handlers:   []HandleFunc{},
		idx:        -1,
	}
}

// Next transfers control to next handler
func (ctx *Context) Next() {
	ctx.idx++
	for ; ctx.idx < len(ctx.handlers); ctx.idx++ {
		ctx.handlers[ctx.idx](ctx)
	}
}

// Query ...
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// PostForm ...
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.FormValue(key)
}

// SetStatus ...
func (ctx *Context) SetStatus(code int) {
	ctx.StatusCode = code
	ctx.Writer.WriteHeader(code)
}

// SetHeader ...
func (ctx *Context) SetHeader(key string, value string) {
	ctx.Writer.Header().Set(key, value)
}

// String ...
func (ctx *Context) String(code int, format string, values ...interface{}) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "text/plain")
	_, _ = ctx.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON ...
func (ctx *Context) JSON(code int, obj interface{}) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(ctx.Writer, err.Error(), 500)
	}
}

// HTML ...
func (ctx *Context) HTML(code int, html string) {
	ctx.SetStatus(code)
	ctx.SetHeader("Content-Type", "text/html")
	_, _ = ctx.Writer.Write([]byte(html))
}

// Data ...
func (ctx *Context) Data(code int, data []byte) {
	ctx.SetStatus(code)
	_, _ = ctx.Writer.Write(data)
}

// Fail ...
func (ctx *Context) Fail(code int, err string) {
	// set idx to the end -> stop handler chain execution when Fail
	ctx.idx = len(ctx.handlers)
	ctx.JSON(code, OBJ{"message": err})
}
