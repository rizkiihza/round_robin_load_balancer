package configs

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	appPort                     string
	pingTimeoutMs               int
	checkPeriodMs               int
	applicationServiceHosts     []string
	applicationServicePingPaths []string
}

func NewStub(
	appPort string,
	pingTimeoutMs int,
	checkPeriodMs int,
	applicationServiceHosts []string,
	applicationServicePingPaths []string,
) *Config {
	return &Config{
		appPort:                     appPort,
		pingTimeoutMs:               pingTimeoutMs,
		checkPeriodMs:               checkPeriodMs,
		applicationServiceHosts:     applicationServiceHosts,
		applicationServicePingPaths: applicationServicePingPaths,
	}
}

func New() *Config {
	appPort := os.Getenv("APP_PORT")
	pingTimeoutMsString := os.Getenv("PING_TIMEOUT_MS")
	checkPeriodMsString := os.Getenv("CHECK_PERIOD_MS")
	applicationServiceHosts := strings.Split(os.Getenv("APPLICATION_SERVICE_HOSTS"), ",")
	applicationServicePingPaths := strings.Split(os.Getenv("APPLICATION_SERVICE_PING_PATHS"), ",")

	pingTimeoutMs, err := strconv.Atoi(pingTimeoutMsString)
	if err != nil {
		panic(err)
	}
	checkPeriodMs, err := strconv.Atoi(checkPeriodMsString)
	if err != nil {
		panic(err)
	}
	return &Config{
		appPort:                     appPort,
		pingTimeoutMs:               pingTimeoutMs,
		checkPeriodMs:               checkPeriodMs,
		applicationServiceHosts:     applicationServiceHosts,
		applicationServicePingPaths: applicationServicePingPaths,
	}
}

func (c *Config) GetAppPort() string {
	return c.appPort
}

func (c *Config) GetPingTimeoutMs() int {
	return c.pingTimeoutMs
}

func (c *Config) GetCheckPeriodMs() int {
	return c.checkPeriodMs
}

func (c *Config) GetApplicationServiceHosts() []string {
	return c.applicationServiceHosts
}

func (c *Config) GetApplicationServicePingPaths() []string {
	return c.applicationServicePingPaths
}
