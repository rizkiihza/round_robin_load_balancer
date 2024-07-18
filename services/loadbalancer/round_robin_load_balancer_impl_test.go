package loadbalancer

import (
	"loadbalancer/model"
	"testing"
	"time"
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
		for i := 0; i < 10000; i++ {
			key, err := lb.GetNextServiceKey(healthyServiceKeys)
			if err != nil {
				t.Errorf("not expecting error when getting next service key")
			}
			result = append(result, key)
		}

		if len(result) < 10000 {
			t.Errorf("expecting 100 keys")
		}
		healthyServiceIndex := 0
		for i := 0; i < 10000; i++ {
			current_i := i
			if result[current_i] != healthyServiceKeys[healthyServiceIndex] {
				t.Errorf("expecting %d result %s to be equal to %s", current_i, result[i], healthyServiceKeys[healthyServiceIndex])
			}
			healthyServiceIndex = (healthyServiceIndex + 1) % int(lb.totalAppNumber)
		}
	})

	t.Run("test get next service key parallel with 30 goroutine", func(t *testing.T) {
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
				"localhost:8002",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-3",
				"localhost:8003",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-4",
				"localhost:8004",
				"/",
				"/",
			),
			model.NewApplicationServiceAddress(
				"app_service-5",
				"localhost:8005",
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
		result_channel := make(chan string)

		// run 30 concurrent goroutine to check
		for i := 0; i < 30; i++ {
			go func(result chan<- string, lb LoadBalancer, healthyServiceKeys []string) {
				for i := 0; i < 100; i++ {
					serviceKey, _ := lb.GetNextServiceKey(healthyServiceKeys)
					result <- serviceKey
				}
			}(result_channel, lb, healthyServiceKeys)
		}

		go func(result_channel chan<- string) {
			time.Sleep(3 * time.Second)
			close(result_channel)
		}(result_channel)
		result_slice := make([]string, 0)

		for val := range result_channel {
			result_slice = append(result_slice, val)
		}

		map_count := make(map[string]int)
		for _, val := range result_slice {
			if _, ok := map_count[val]; !ok {
				map_count[val] = 0
			}
			map_count[val] += 1
		}

		min_value := -1
		max_value := -1

		for _, v := range map_count {
			if min_value == -1 || v < min_value {
				min_value = v
			}
			if max_value == -1 || v > max_value {
				max_value = v
			}
		}

		if max_value-min_value > 1 {
			t.Errorf("expect the difference of most used and least used to be less than equal to 1")
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
