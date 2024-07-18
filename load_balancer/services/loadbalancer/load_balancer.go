package loadbalancer

type LoadBalancer interface {
	GetNextServiceKey(healthyServiceKeys []string) (string, error)
}
