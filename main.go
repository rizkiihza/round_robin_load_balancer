package main

import (
	"fmt"
	"loadbalancer/clients"
	"loadbalancer/configs"
	"loadbalancer/handler"
	"loadbalancer/model"
	"loadbalancer/services/healthcheck"
	"loadbalancer/services/loadbalancer"
	"loadbalancer/services/processor"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := configs.New()
	appClient := clients.New()
	appAddresses := createApplicationServiceAddresses(config)
	healthcheck := healthcheck.New(
		appAddresses,
		appClient,
		config.GetPingTimeoutMs(),
		config.GetCheckPeriodMs())
	loadbalancer := loadbalancer.New(
		appAddresses)
	processor := processor.New(
		loadbalancer,
		healthcheck,
		createApplicationAddressesMap(appAddresses),
		appClient)

	// run periodic healthcheck
	done := make(chan struct{})
	healthcheck.PeriodicalCheck(done)

	// register http handler
	handler := handler.New(processor)
	http.HandleFunc("/", handler.Post)

	// graceful shutdown
	gracefulShutdown(done)

	// listen to http request
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetAppPort()), nil)
	if err != nil {
		panic(err)
	}
}

func gracefulShutdown(done chan<- struct{}) {
	go func(done chan<- struct{}) {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		done <- struct{}{}
	}(done)
}

func createApplicationServiceAddresses(config *configs.Config) []*model.ApplicationServiceAddress {
	addresses := make([]*model.ApplicationServiceAddress, 0)
	for i := 0; i < len(config.GetApplicationServiceHosts()); i++ {
		address := model.NewApplicationServiceAddress(
			fmt.Sprintf("app_service-%d", i),
			config.GetApplicationServiceHosts()[i],
			config.GetApplicationServiceCallPaths()[i],
			config.GetApplicationServicePingPaths()[i])
		addresses = append(addresses, address)
	}
	return addresses
}

func createApplicationAddressesMap(appAddresses []*model.ApplicationServiceAddress) map[string]*model.ApplicationServiceAddress {
	appAddressesMap := make(map[string]*model.ApplicationServiceAddress)
	for i := 0; i < len(appAddresses); i++ {
		appAddress := appAddresses[i]
		appAddressesMap[appAddress.GetKey()] = appAddress
	}
	return appAddressesMap
}
