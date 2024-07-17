package clients

import (
	"context"
	"fmt"
	"loadbalancer/model"
	"net/http"
)

type ApplicationServiceImpl struct {
}

func (as *ApplicationServiceImpl) Call(ctx context.Context, request *http.Request, appService *model.ApplicationServiceAddress) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", appService.GetHost(), appService.GetCallPath())
	req, err := http.NewRequestWithContext(ctx, "GET", url, request.Body)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (as *ApplicationServiceImpl) Ping(ctx context.Context, appService *model.ApplicationServiceAddress) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", appService.GetHost(), appService.GetCallPath())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
