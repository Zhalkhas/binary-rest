package http

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/zhalkhas/binary-rest/internal/repos"
	"log/slog"
	"net/http"
	"strconv"
)

const endpointValuePathParam = "value"

type IndicesController struct {
	handler  http.Handler
	indxRepo repos.Indices
}

func NewIndicesController(indices repos.Indices, middlewares ...func(http.Handler) http.Handler) IndicesController {
	m := chi.NewRouter()
	m.Use(middlewares...)
	ir := IndicesController{m, indices}
	m.Get(fmt.Sprintf("/endpoint/{%s}", endpointValuePathParam), ir.getIndexHandler)
	return ir
}

func (i IndicesController) getIndexHandler(rw http.ResponseWriter, r *http.Request) {
	value, err := strconv.ParseInt(chi.URLParam(r, endpointValuePathParam), 10, 64)
	if err != nil {
		if err = render.Render(rw, r, ErrInvalidValue); err != nil {
			slog.Error("error while rendering invalid value error", "err", err)
			return
		}
		slog.Error(
			"error while parsing value from path",
			"err", err, "value", chi.URLParam(r, endpointValuePathParam),
		)
		return
	}
	index, err := i.indxRepo.Search(r.Context(), value)
	if err == nil {
		if err = render.Render(rw, r, IndexFoundResponse{Index: index, Value: value}); err != nil {
			slog.Error(
				"error while rendering index found response",
				"err", err, "index", index, "value", value,
			)
			return
		}
		return
	} else if errors.Is(err, repos.ErrIndexNotFound) {
		if err := render.Render(rw, r, ErrIndexNotFound); err != nil {
			slog.Error(
				"error while rendering index not found error",
				"err", err, "value", value,
			)
			return
		}
		slog.Error(
			"index not found error",
			"err", err, "value", value,
		)
		return
	} else {
		if err = render.Render(rw, r, ErrUnknown); err != nil {
			slog.Error(
				"error while rendering unknown error",
				"err", err, "value", value,
			)
			return
		}
		slog.Error(
			"unknown error happened during search",
			"err", err, "value", value,
		)
		return
	}
}

func (i IndicesController) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	i.handler.ServeHTTP(rw, r)
}
