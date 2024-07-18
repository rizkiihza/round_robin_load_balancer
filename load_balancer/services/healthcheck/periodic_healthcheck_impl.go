package healthcheck

import (
	"context"
	"loadbalancer/clients"
	"loadbalancer/model"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var LOGGER = log.New(os.Stdout, "[HealthCheck] INFO:", log.Ldate)

type PeriodicHealthcheck struct {
	appAddresses          []*model.ApplicationServiceAddress
	appClient             clients.ApplicationServiceClient
	healthyAppServiceKeys []string
	pingTimeoutMs         int
	checkPeriodMs         int
	rwMutex               *sync.RWMutex
}

func New(
	appAddresses []*model.ApplicationServiceAddress,
	appClient clients.ApplicationServiceClient,
	pingTimeoutMs int,
	checkPeriodMs int) *PeriodicHealthcheck {

	return &PeriodicHealthcheck{
		appAddresses:          appAddresses,
		appClient:             appClient,
		healthyAppServiceKeys: make([]string, 0),
		pingTimeoutMs:         pingTimeoutMs,
		checkPeriodMs:         checkPeriodMs,
		rwMutex:               &sync.RWMutex{},
	}
}

func (p *PeriodicHealthcheck) PeriodicalCheck(done <-chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				LOGGER.Println("got done signal")
				return
			case <-time.After(time.Duration(p.checkPeriodMs) * time.Millisecond):
				p.check(context.Background())
				if len(p.healthyAppServiceKeys) > 0 {
					LOGGER.Printf("done doing periodic health check, healthy services: %s", strings.Join(p.healthyAppServiceKeys, ","))
				} else {
					LOGGER.Println("done doing checking: no healthy service")
				}
			}
		}
	}()
}

func (p *PeriodicHealthcheck) check(ctx context.Context) {
	healthyServiceKeys := make([]string, 0)
	for _, address := range p.appAddresses {
		each_ctx, cancel := context.WithTimeout(ctx, time.Duration(p.pingTimeoutMs)*time.Millisecond)
		defer cancel()

		response, err := p.appClient.Ping(each_ctx, address)
		if err != nil {
			continue
		}

		if response.StatusCode != http.StatusOK {
			continue
		}
		healthyServiceKeys = append(healthyServiceKeys, address.GetKey())
	}

	// use RW mutex to prevent race condition
	p.rwMutex.Lock()
	p.healthyAppServiceKeys = healthyServiceKeys
	p.rwMutex.Unlock()
}

func (p *PeriodicHealthcheck) GetHealthyServiceKeys() []string {
	// use RW mutex to prevent race condition
	p.rwMutex.RLock()
	result := p.healthyAppServiceKeys
	p.rwMutex.RUnlock()
	return result
}
