package loadbalancer

import (
	"fmt"
	"loadbalancer/model"
	"sync/atomic"
)

type RoundRobinLoadBalancerImpl struct {
	appAddress     []*model.ApplicationServiceAddress
	totalAppNumber int64
	index          int64
}

func New(appAddress []*model.ApplicationServiceAddress) *RoundRobinLoadBalancerImpl {
	return &RoundRobinLoadBalancerImpl{
		appAddress:     appAddress,
		totalAppNumber: int64(len(appAddress)),
		index:          -1,
	}
}

func (lb *RoundRobinLoadBalancerImpl) GetNextServiceKey(healthyServiceKeys []string) (string, error) {
	for {
		if len(healthyServiceKeys) == 0 {
			return "", fmt.Errorf("no healthy application services")
		}
		next_key := lb.getNextKey()
		if contains(healthyServiceKeys, next_key) {
			return next_key, nil
		}
	}
}

func (lb *RoundRobinLoadBalancerImpl) getNextKey() string {
	// add one to atomic int64
	value := atomic.AddInt64(&lb.index, 1)

	// if value more than a hundred times total app number, reduce by total
	// we do this to have a throttling on max value of the atomic int64
	if value > 100*lb.totalAppNumber {
		atomic.AddInt64(&lb.index, -lb.totalAppNumber)
	}

	// do reconciliation on key value
	key := value % lb.totalAppNumber
	if key < 0 {
		key = key + lb.totalAppNumber
	}
	return lb.appAddress[key].GetKey()
}

func contains(serviceKeys []string, key string) bool {
	for _, skey := range serviceKeys {
		if skey == key {
			return true
		}
	}
	return false
}
