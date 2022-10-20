// Package appconfig implements a main routine that is the starting point of
// location-mgmt microservice.
package appconfig

import (
	"context"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	middleKit "github.com/oboadagd/kit-go/middleware/echo"
	"github.com/oboadagd/location-mgmt/controller"
	"github.com/oboadagd/location-mgmt/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartApp implements configuration and start-up of microservice.
func StartApp() {

	echoInstance := echo.New()

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(echoInstance)

	echoInstance.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency_human=${latency_human}\n",
	}))
	echoInstance.Use(middleware.Recover())

	locationController := controller.NewLocationController()

	errorHandlerMiddle := middleKit.NewErrorHandlerMiddleware()

	r := router.NewRouter(echoInstance, locationController, errorHandlerMiddle)
	r.Init()

	// Start server
	go func() {
		if err := echoInstance.Start(":8081"); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGSTOP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := echoInstance.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
