package http

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

type IndicesRouter struct {
	*chi.Mux
}

func NewIndicesRouter() IndicesRouter {
	m := chi.NewRouter()
	ir := IndicesRouter{m}
	m.Get("/endpoint/{index}", ir.getIndexHandler)
	return ir
}

func (i IndicesRouter) getIndexHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func (i IndicesRouter) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	i.Mux.ServeHTTP(rw, r)
}
