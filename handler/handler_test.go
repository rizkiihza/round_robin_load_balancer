package handler

import (
	"bytes"
	"fmt"
	"io"
	"loadbalancer/services/processor/mock_processor"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHandler_Post(t *testing.T) {
	t.Run("500 response from application service", func(t *testing.T) {
		processor := mock_processor.NewMockProcessor(gomock.NewController(t))

		processor.EXPECT().ForwardRequest(gomock.Any(), gomock.Any()).AnyTimes().
			Return(nil, fmt.Errorf("got error from upstream"))
		h := &Handler{
			processor: processor,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "localhost:8000", nil)
		h.Post(w, r)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Errorf("not expecting error when getting data")
		}

		expectedResponse := "500 - Internal Server error"
		if string(data[:]) != expectedResponse {
			t.Errorf("expecting message to be %s", expectedResponse)
		}
	})

	t.Run("Success response from application service", func(t *testing.T) {
		processor := mock_processor.NewMockProcessor(gomock.NewController(t))

		processor_response := http.Response{
			Body: io.NopCloser(bytes.NewBufferString("Hello World")),
		}
		processor.EXPECT().ForwardRequest(gomock.Any(), gomock.Any()).AnyTimes().
			Return(&processor_response, nil)
		h := &Handler{
			processor: processor,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "localhost:8000", nil)
		h.Post(w, r)
		res := w.Result()
		data, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			t.Errorf("not expecting error when getting data")
		}

		expectedResponse := "Hello World"
		if string(data[:]) != expectedResponse {
			t.Errorf("expecting message to be %s", expectedResponse)
		}
	})

}
