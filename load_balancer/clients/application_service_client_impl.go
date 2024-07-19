package clients

import (
	"context"
	"fmt"
	"loadbalancer/model"
	"log"
	"net/http"
	"os"
)

var LOGGER = log.New(os.Stdout, "client INFO:", log.Ldate)

type ApplicationServiceClientImpl struct {
}

func New() *ApplicationServiceClientImpl {
	return &ApplicationServiceClientImpl{}
}

func (as *ApplicationServiceClientImpl) Call(ctx context.Context, request *http.Request, appService *model.ApplicationServiceAddress) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", appService.GetHost(), appService.GetCallPath())
	req, err := http.NewRequestWithContext(ctx, "POST", url, request.Body)
	if err != nil {
		LOGGER.Println("got error when creating request", err.Error())
		return nil, err
	}

	LOGGER.Printf("http request call to %s", url)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		LOGGER.Println("got error when doing http request", err.Error())
		return nil, err
	}
	return response, nil
}

func (as *ApplicationServiceClientImpl) Ping(ctx context.Context, appService *model.ApplicationServiceAddress) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", appService.GetHost(), appService.GetPingPath())
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
