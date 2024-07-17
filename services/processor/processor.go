package processor

import (
	"context"
	"net/http"
)

type Processor interface {
	ForwardRequest(ctx context.Context, request *http.Request) (*http.Response, error)
}
