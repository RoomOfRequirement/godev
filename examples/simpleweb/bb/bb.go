package bb

import (
	"net/http"
	"strings"
)

// HandleFunc ...
type HandleFunc func(*Context)

// OBJ ...
type OBJ map[string]interface{}

// Engine ...
type Engine struct {
	router *router

	*RouterGroup                // Engine as root router group
	groups       []*RouterGroup // store all groups
}

// Default use Logger() and Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// New ...
func New() *Engine {
	engine := &Engine{router: newRouter()}
	// root
	engine.RouterGroup = &RouterGroup{
		prefix:      "",
		middlewares: nil,
		parent:      nil,
		engine:      engine,
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) addRoute(method, path string, handler HandleFunc) {
	e.router.addRoute(method, path, handler)
}

// GET method
func (e *Engine) GET(path string, handler HandleFunc) {
	e.addRoute("GET", path, handler)
}

// POST method
func (e *Engine) POST(path string, handler HandleFunc) {
	e.addRoute("POST", path, handler)
}

// Run server
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP ...
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandleFunc
	for _, gr := range e.groups {
		// dispatch request to corresponding group
		// and add corresponding middlewares to context
		if strings.HasPrefix(r.URL.Path, gr.prefix) {
			middlewares = append(middlewares, gr.middlewares...)
		}
	}
	ctx := newContext(w, r)
	ctx.handlers = middlewares

	// handle request
	e.router.handle(ctx)
}

// RouterGroup ...
type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc // middleware is binding on group
	parent      *RouterGroup // for nesting
	engine      *Engine      // one engine for all group
}

// AddGroup adds new router group with prefix
func (rg *RouterGroup) AddGroup(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix:      rg.prefix + prefix,
		middlewares: []HandleFunc{},
		parent:      rg,
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) addRoute(method, relatedPath string, handler HandleFunc) {
	rg.engine.router.addRoute(method, rg.prefix+relatedPath, handler)
}

// GET method
func (rg *RouterGroup) GET(path string, handler HandleFunc) {
	rg.addRoute("GET", path, handler)
}

// POST method
func (rg *RouterGroup) POST(path string, handler HandleFunc) {
	rg.addRoute("POST", path, handler)
}

// Use sets middlewares to router group
func (rg *RouterGroup) Use(middlewares ...HandleFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
}
