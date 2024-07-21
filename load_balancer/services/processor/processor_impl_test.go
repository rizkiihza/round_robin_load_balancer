package processor

import (
	"bytes"
	"context"
	"io"
	mock_clients "loadbalancer/clients/mock"
	"loadbalancer/model"
	mock_healthcheck "loadbalancer/services/healthcheck/mock"
	mock_loadbalancer "loadbalancer/services/loadbalancer/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestProcessorImpl_ForwardRequest(t *testing.T) {
	t.Run("forward request, all dependency behave well", func(t *testing.T) {
		addresses := []*model.ApplicationServiceAddress{
			model.NewApplicationServiceAddress(
				"app_service-0",
				"localhost:8000",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-1",
				"localhost:8001",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-2",
				"localhost:8001",
				"/",
			),
		}

		addressMap := make(map[string]*model.ApplicationServiceAddress)
		for _, address := range addresses {
			addressMap[address.GetKey()] = address
		}

		gomockCtrl := gomock.NewController(t)
		lb := mock_loadbalancer.NewMockLoadBalancer(gomockCtrl)
		lb.EXPECT().GetNextServiceKey(gomock.Any()).AnyTimes().
			Return("app_service-0", nil)
		healthcheck := mock_healthcheck.NewMockHealthcheck(gomockCtrl)
		healthcheck.EXPECT().GetHealthyServiceKeys().AnyTimes().
			Return([]string{
				"app_service-0",
				"app_service-1",
				"app_service-2",
			})
		processor_response := http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("Hello World")),
		}
		clients := mock_clients.NewMockApplicationServiceClient(gomockCtrl)
		clients.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
			Return(&processor_response, nil)
		p := ProcessorImpl{
			lb:           lb,
			healthcheck:  healthcheck,
			appAddresses: addressMap,
			appClient:    clients,
		}

		request := httptest.NewRequest("POST", "/call",
			io.NopCloser(bytes.NewBufferString("Hello World")))
		res, err := p.ForwardRequest(context.Background(), request)
		if err != nil {
			t.Errorf("expecting error to be nil")
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("expecting success status code, instead getting %d", res.StatusCode)
		}
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

	t.Run("forward request, no healthy services", func(t *testing.T) {
		addresses := []*model.ApplicationServiceAddress{
			model.NewApplicationServiceAddress(
				"app_service-0",
				"localhost:8000",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-1",
				"localhost:8001",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-2",
				"localhost:8001",
				"/",
			),
		}

		addressMap := make(map[string]*model.ApplicationServiceAddress)
		for _, address := range addresses {
			addressMap[address.GetKey()] = address
		}

		gomockCtrl := gomock.NewController(t)
		lb := mock_loadbalancer.NewMockLoadBalancer(gomockCtrl)
		lb.EXPECT().GetNextServiceKey(gomock.Any()).AnyTimes().
			Return("app_service-0", nil)
		healthcheck := mock_healthcheck.NewMockHealthcheck(gomockCtrl)
		healthcheck.EXPECT().GetHealthyServiceKeys().AnyTimes().
			Return([]string{})
		processor_response := http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("Hello World")),
		}
		clients := mock_clients.NewMockApplicationServiceClient(gomockCtrl)
		clients.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
			Return(&processor_response, nil)
		p := ProcessorImpl{
			lb:           lb,
			healthcheck:  healthcheck,
			appAddresses: addressMap,
			appClient:    clients,
		}

		request := httptest.NewRequest("POST", "/call",
			io.NopCloser(bytes.NewBufferString("Hello World")))
		_, err := p.ForwardRequest(context.Background(), request)
		if err == nil {
			t.Errorf("expecting error to be nil")
		}
		expectedErrorMessage := "there is no healthy services"
		if err.Error() != expectedErrorMessage {
			t.Errorf("expecting error message to be %s, instead getting %s", expectedErrorMessage, err.Error())
		}
	})
}
