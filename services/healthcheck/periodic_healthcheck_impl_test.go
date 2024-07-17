package healthcheck

import (
	"bytes"
	"io"
	mock_clients "loadbalancer/clients/mock"
	"loadbalancer/model"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestPeriodicHealthcheck_PeriodicalCheck(t *testing.T) {
	t.Run("test periodic health check, services are alive", func(t *testing.T) {
		appClient := mock_clients.NewMockApplicationServiceClient(gomock.NewController(t))
		mock_response := http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("ping success")),
		}
		appClient.EXPECT().Ping(gomock.Any(), gomock.Any()).
			AnyTimes().Return(&mock_response, nil)
		done := make(chan struct{})
		addresses := []*model.ApplicationServiceAddress{
			model.NewApplicationServiceAddress(
				"app_service-0",
				"localhost:8000",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-1",
				"localhost:8001",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-2",
				"localhost:8001",
				"/",
				"/",
			),
		}

		p := &PeriodicHealthcheck{
			appAddresses:          addresses,
			appClient:             appClient,
			healthyAppServiceKeys: make([]string, 0),
			pingTimeoutMs:         100,
			checkPeriodMs:         100,
			rwMutex:               &sync.RWMutex{},
		}
		p.PeriodicalCheck(done)
		time.Sleep(2 * time.Second)
		if len(p.healthyAppServiceKeys) != 3 {
			t.Errorf("expecting 3 healthy service")
		}
		if p.healthyAppServiceKeys[0] != "app_service-0" {
			t.Errorf("expecing first key to be app-service-0")
		}
		if p.healthyAppServiceKeys[1] != "app_service-1" {
			t.Errorf("expecing second key to be app-service-1")
		}
		if p.healthyAppServiceKeys[2] != "app_service-2" {
			t.Errorf("expecing third key to be app-service-2")
		}
		done <- struct{}{}
	})

	t.Run("test periodic health check, services are dying", func(t *testing.T) {
		appClient := mock_clients.NewMockApplicationServiceClient(gomock.NewController(t))
		mock_response := http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("500 - Internal Server Error")),
		}
		appClient.EXPECT().Ping(gomock.Any(), gomock.Any()).
			AnyTimes().Return(&mock_response, nil)
		done := make(chan struct{})
		addresses := []*model.ApplicationServiceAddress{
			model.NewApplicationServiceAddress(
				"app_service-0",
				"localhost:8000",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-1",
				"localhost:8001",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-1",
				"localhost:8001",
				"/",
				"/",
			),
		}

		p := &PeriodicHealthcheck{
			appAddresses:          addresses,
			appClient:             appClient,
			healthyAppServiceKeys: make([]string, 0),
			pingTimeoutMs:         100,
			checkPeriodMs:         100,
			rwMutex:               &sync.RWMutex{},
		}
		p.PeriodicalCheck(done)
		time.Sleep(2 * time.Second)
		if len(p.healthyAppServiceKeys) != 0 {
			t.Errorf("expecting 0 healthy service")
		}
		done <- struct{}{}
	})

}

func TestPeriodicHealthcheck_GetHealthyServiceKeys(t *testing.T) {
	t.Run("Get healthy service keys", func(t *testing.T) {
		p := &PeriodicHealthcheck{
			healthyAppServiceKeys: []string{"app_service-1", "app_service-2", "app_service-3"},
			pingTimeoutMs:         100,
			checkPeriodMs:         100,
			rwMutex:               &sync.RWMutex{},
		}

		result := p.GetHealthyServiceKeys()
		if len(result) != 3 {
			t.Errorf("expecting 3 healthy services")
		}
		if p.healthyAppServiceKeys[0] != "app_service-1" {
			t.Errorf("expecing first key to be app-service-1")
		}
		if p.healthyAppServiceKeys[1] != "app_service-2" {
			t.Errorf("expecing second key to be app-service-2")
		}
		if p.healthyAppServiceKeys[2] != "app_service-3" {
			t.Errorf("expecing third key to be app-service-3")
		}
	})
}
