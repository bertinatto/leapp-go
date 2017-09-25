package web

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/leapp-to/leapp-go/pkg/executor"
	"github.com/leapp-to/leapp-go/pkg/msg"
)

type Handler struct {
	mux     *http.ServeMux
	context context.Context
	errorCh chan error
	//options Options
}

func (h *Handler) Run() {
	srv := &http.Server{
		Addr:    ":8000",
		Handler: h.mux,
	}

	if listener, err := net.Listen("tcp", ":8000"); err == nil {
		h.errorCh <- srv.Serve(listener)
	} else {
		h.errorCh <- err
	}
}

func (h *Handler) ErrorCh() <-chan error {
	return h.errorCh
}

func New() *Handler {
	h := &Handler{
		mux:     http.NewServeMux(),
		errorCh: make(chan error),
	}

	h.mux.HandleFunc("/migrate-machine", func(w http.ResponseWriter, req *http.Request) {
		m := msg.MigrateMachine{}
		json.NewDecoder(req.Body).Decode(&m)

		r := executor.Run(&m)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(r)
	})

	return h
}
