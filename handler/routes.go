package handler

import (
	"github.com/cnosuke/imagine/entity"
	"github.com/gorilla/mux"
	"net/http"
)

type route struct {
	method  string
	path    string
	handler func(http.ResponseWriter, *http.Request) (int, interface{}, error)
}

const APIPrefix = "/v1"

func (h *Handler) Routing() *mux.Router {
	// Health check
	h.router.
		Methods(http.MethodGet).
		Path(APIPrefix + "/health").
		Handler(h.getChain().Then(_handler{h.healthcheckHandler}))

	return h.router
}

func (h *Handler) healthcheckHandler(_ *http.Request) (int, interface{}, error) {
	ctx := h.context
	revision := ctx.Value("revision").(string)

	return http.StatusOK, entity.Health{Revision: revision}, nil
}
