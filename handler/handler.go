package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	router *mux.Router
	context context.Context
	chain        alice.Chain
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	PathPrefix = "/api"
)

func NewHandler(ctx context.Context) *Handler {
	h := &Handler{
		context: ctx,
		router: mux.NewRouter().PathPrefix(PathPrefix).Subrouter(),
	}

	h.setChain(
		alice.New(
			appLoggingMiddleware(zap.L()),
		),
	)

	return h
}

func (h *Handler) setChain(chain alice.Chain) {
	h.chain = chain
	return
}

func (h *Handler) getChain() alice.Chain {
	return h.chain
}

func appLoggingMiddleware(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tStart := time.Now()
			next.ServeHTTP(w, r)
			tEnd := time.Now()
			t := tEnd.Sub(tStart)
			l.Sugar().Infof("[%s] %s %s", r.Method, r.URL, t.String())
		})
	}
}

type _handler struct {
	h func(*http.Request) (int, interface{}, error)
}

func (h _handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	status, res, err := h.h(r)
	if err != nil {
		zap.S().Infof("error: %s", err)
		w.WriteHeader(status)
		encoder.Encode(res)
		return
	}
	w.WriteHeader(status)

	if res != nil {
		encoder.Encode(res)
	}
	return
}
