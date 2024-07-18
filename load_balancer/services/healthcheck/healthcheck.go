package healthcheck

type Healthcheck interface {
	GetHealthyServiceKeys() []string
	PeriodicalCheck(done <-chan struct{})
}
