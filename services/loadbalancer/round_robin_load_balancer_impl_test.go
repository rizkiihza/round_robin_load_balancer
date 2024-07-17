package loadbalancer

import (
	"loadbalancer/model"
	"testing"
)

func TestRoundRobinLoadBalancerImpl_GetNextServiceKey(t *testing.T) {
	t.Run("test get next service key sequential", func(t *testing.T) {
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

		lb := &RoundRobinLoadBalancerImpl{
			appAddress:     addresses,
			totalAppNumber: 3,
			index:          -1,
		}
		healthyServiceKeys := []string{"app_service-0", "app_service-1", "app_service-2"}

		result := make([]string, 0)
		for i := 0; i < 1000; i++ {
			key, err := lb.GetNextServiceKey(healthyServiceKeys)
			if err != nil {
				t.Errorf("not expecting error when getting next service key")
			}
			result = append(result, key)
		}

		if len(result) < 1000 {
			t.Errorf("expecting 100 keys")
		}
		healthyServiceIndex := 0
		for i := 0; i < 1000; i++ {
			current_i := i
			if result[current_i] != healthyServiceKeys[healthyServiceIndex] {
				t.Errorf("expecting %d result %s to be equal to %s", current_i, result[i], healthyServiceKeys[healthyServiceIndex])
			}
			healthyServiceIndex = (healthyServiceIndex + 1) % int(lb.totalAppNumber)
		}
	})

	t.Run("test get next service key, no healthy service exist", func(t *testing.T) {
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

		lb := &RoundRobinLoadBalancerImpl{
			appAddress:     addresses,
			totalAppNumber: 3,
			index:          -1,
		}
		healthyServiceKeys := []string{}

		_, err := lb.GetNextServiceKey(healthyServiceKeys)
		if err == nil {
			t.Errorf("expecting error here, since there is no healthy service")
		}

		expectedErrorMessage := "no healthy application services"
		if err.Error() != expectedErrorMessage {
			t.Errorf("expecting error with message: %s", expectedErrorMessage)
		}
	})

	t.Run("test get next service key sequential, some service are down", func(t *testing.T) {
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
			model.NewApplicationServiceAddress(
				"app_service-3",
				"localhost:8001",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-4",
				"localhost:8001",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-5",
				"localhost:8001",
				"/",
				"/",
			),
		}

		lb := &RoundRobinLoadBalancerImpl{
			appAddress:     addresses,
			totalAppNumber: int64(len(addresses)),
			index:          -1,
		}
		healthyServiceKeys := []string{"app_service-0", "app_service-3", "app_service-5"}

		result := make([]string, 0)
		for i := 0; i < 1000; i++ {
			key, err := lb.GetNextServiceKey(healthyServiceKeys)
			if err != nil {
				t.Errorf("not expecting error when getting next service key")
			}
			result = append(result, key)
		}

		if len(result) < 1000 {
			t.Errorf("expecting 10 keys")
		}
		healthyServiceIndex := 0
		for i := 0; i < 1000; i++ {
			current_i := i
			if result[current_i] != healthyServiceKeys[healthyServiceIndex] {
				t.Errorf("expecting %d result %s to be equal to %s", current_i, result[i], healthyServiceKeys[healthyServiceIndex])
			}
			healthyServiceIndex = (healthyServiceIndex + 1) % len(healthyServiceKeys)
		}
	})
}
