package handler

import (
	"io"
	"loadbalancer/services/processor"
	"net/http"
)

type Handler struct {
	processor processor.Processor
}

func New(processor processor.Processor) *Handler {
	return &Handler{
		processor: processor,
	}
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	response, err := h.processor.ForwardRequest(r.Context(), r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 - Internal Server error"))
	} else {
		io.Copy(w, response.Body)
		response.Body.Close()
	}
}
