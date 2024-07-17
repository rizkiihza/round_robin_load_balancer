package healthcheck

type Healthcheck interface {
	GetHealthyServiceKeys() []string
}
