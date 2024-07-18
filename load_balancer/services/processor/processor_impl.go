package processor

import (
	"context"
	"fmt"
	"loadbalancer/clients"
	"loadbalancer/model"
	"loadbalancer/services/healthcheck"
	"loadbalancer/services/loadbalancer"
	"net/http"
)

type ProcessorImpl struct {
	lb           loadbalancer.LoadBalancer
	healthcheck  healthcheck.Healthcheck
	appAddresses map[string]*model.ApplicationServiceAddress
	appClient    clients.ApplicationServiceClient
}

func New(
	lb loadbalancer.LoadBalancer,
	healthcheck healthcheck.Healthcheck,
	appAddresses map[string]*model.ApplicationServiceAddress,
	appClient clients.ApplicationServiceClient) *ProcessorImpl {

	return &ProcessorImpl{
		lb:           lb,
		healthcheck:  healthcheck,
		appAddresses: appAddresses,
		appClient:    appClient,
	}
}

func (p ProcessorImpl) ForwardRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	// get list of healthy service
	healthy_service_keys := p.healthcheck.GetHealthyServiceKeys()
	if len(healthy_service_keys) == 0 {
		return nil, fmt.Errorf("there is no healthy services")
	}

	// get next service to call
	next_service_key, err := p.lb.GetNextServiceKey(healthy_service_keys)
	if err != nil {
		return nil, err
	}

	// get address and call the application service
	appAddress := p.appAddresses[next_service_key]
	response, err := p.appClient.Call(ctx, request, appAddress)
	if err != nil {
		return nil, err
	}

	return response, nil
}
