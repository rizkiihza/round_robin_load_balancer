package model

type ApplicationServiceAddress struct {
	key      string
	host     string
	callPath string
	pingPath string
}

func New(key string, host string, callPath string, pingPath string) *ApplicationServiceAddress {
	return &ApplicationServiceAddress{
		key:      key,
		host:     host,
		callPath: callPath,
		pingPath: pingPath,
	}
}

func (as ApplicationServiceAddress) GetKey() string {
	return as.key
}

func (as ApplicationServiceAddress) GetHost() string {
	return as.host
}

func (as ApplicationServiceAddress) GetCallPath() string {
	return as.callPath
}

func (as ApplicationServiceAddress) GetPingPath() string {
	return as.pingPath
}
