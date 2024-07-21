package model

type ApplicationServiceAddress struct {
	key      string
	host     string
	pingPath string
}

func NewApplicationServiceAddress(key string, host string, pingPath string) *ApplicationServiceAddress {
	return &ApplicationServiceAddress{
		key:      key,
		host:     host,
		pingPath: pingPath,
	}
}

func (as ApplicationServiceAddress) GetKey() string {
	return as.key
}

func (as ApplicationServiceAddress) GetHost() string {
	return as.host
}

func (as ApplicationServiceAddress) GetPingPath() string {
	return as.pingPath
}
