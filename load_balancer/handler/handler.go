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

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	response, err := h.processor.ForwardRequest(r.Context(), r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 - Internal Server error\n"))
		return
	}
	defer response.Body.Close()

	io.Copy(w, response.Body)
}
