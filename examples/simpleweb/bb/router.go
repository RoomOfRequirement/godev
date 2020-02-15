package bb

import (
	"godev/basic/datastructure/tree/trietree"
	"net/http"
)

type router struct {
	// trie
	trie trietree.Trie
}

func newRouter() *router {
	return &router{
		trie: trietree.New('/'),
	}
}

func (r *router) addRoute(method, path string, handler HandleFunc) {
	key := method + "-" + path
	r.trie.Put(key, handler)
}

func (r *router) getRoute(method, path string) (handler HandleFunc) {
	key := method + "-" + path
	h := r.trie.Get(key)
	if h == nil {
		return nil
	}
	return h.(HandleFunc)
}

func (r *router) handle(ctx *Context) {
	if handler := r.getRoute(ctx.Method, ctx.Path); handler != nil {
		ctx.handlers = append(ctx.handlers, handler)
	} else {
		ctx.handlers = append(ctx.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", ctx.Path)
		})
	}
	// after executing current handler, pass control to next
	ctx.Next()
}
