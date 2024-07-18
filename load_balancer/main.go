package main

import (
	"context"
	"errors"
	"fmt"
	"loadbalancer/clients"
	"loadbalancer/configs"
	"loadbalancer/handler"
	"loadbalancer/model"
	"loadbalancer/services/healthcheck"
	"loadbalancer/services/loadbalancer"
	"loadbalancer/services/processor"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var LOGGER = log.New(os.Stdout, "INFO:", log.Ldate)

func main() {
	config := configs.NewStub(
		"8086",
		100,
		100,
		[]string{"http://localhost:8081", "http://localhost:8082", "http://localhost:8083", "http://localhost:8084"},
		[]string{"/call", "/call", "/call", "/call"},
		[]string{"/ping", "/ping", "/ping", "/ping"},
	)
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
	LOGGER.Println("starting periodic healthcheck")
	done := make(chan struct{})
	healthcheck.PeriodicalCheck(done)

	// register http handler
	LOGGER.Println("register http handler")
	httpServer := http.Server{
		Addr: fmt.Sprintf(":%s", config.GetAppPort()),
	}

	handler := handler.New(processor)
	http.HandleFunc("/", handler.HandleRequest)

	LOGGER.Println("register graceful shutdown")
	// graceful shutdown
	gracefulShutdown(done, &httpServer)

	// listen to http request
	LOGGER.Println("listen and serve http request on port: ", config.GetAppPort())
	err := httpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		LOGGER.Println("http server has been shutdown")
	} else if err != nil {
		panic(err)
	}
}

func gracefulShutdown(done chan<- struct{}, httpServer *http.Server) {
	go func(done chan<- struct{}) {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		LOGGER.Println("got signal to kill process")
		done <- struct{}{}
		close(done)
		LOGGER.Println("closing done channel")

		LOGGER.Println("http server shutdown")
		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}
		log.Println("Graceful shutdown complete.")
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
