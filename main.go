/*
This is an example Golang program for getting a specific Stock data

Joshua Snyder 10-03-2022

*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cheetahfox/microservice-test/config"
	"github.com/cheetahfox/microservice-test/health"
	"github.com/cheetahfox/microservice-test/router"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := config.Startup()

	mst := fiber.New(config.FiberConfig)

	router.SetupRoutes(mst)

	go func() {
		health.ApiReady = true
		mst.Listen(":2200")
	}()

	// Listen for Sigint or SigTerm and exit if you get them.
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	fmt.Println("Shutdown Started")
	mst.Shutdown()
}
