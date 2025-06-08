package server

import (
	"net/http"
)

type router struct {
	mux        *http.ServeMux
	middleware []Middleware
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func NewRouter() *router {
	return &router{
		mux:        http.NewServeMux(),
		middleware: make([]Middleware, 0),
	}
}

func (r *router) Use(mw ...Middleware) {
	r.middleware = append(r.middleware, mw...)
}

func (r *router) HandleFunc(pattern string, handler http.HandlerFunc) {
	final := handler

	for i := len(r.middleware) - 1; i >= 0; i-- {
		mw := r.middleware[i]
		final = mw(final)
	}
	r.mux.HandleFunc(pattern, final)
}

func (r *router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
