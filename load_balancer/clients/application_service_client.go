package clients

import (
	"context"
	"loadbalancer/model"
	"net/http"
)

type ApplicationServiceClient interface {
	Call(ctx context.Context, request *http.Request, appService *model.ApplicationServiceAddress) (*http.Response, error)
	Ping(ctx context.Context, appService *model.ApplicationServiceAddress) (*http.Response, error)
}
