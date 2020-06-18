package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreasM009/eventstore/pkg/factories"

	"github.com/AndreasM009/eventstore/pkg/injector"
)

func main() {
	// first of all, create the Kuberenetes client
	kubeClient := factories.CreateKubeClient()

	ctx, cancel := context.WithCancel(context.Background())
	sigchannel := make(chan os.Signal)
	signal.Notify(sigchannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	i := injector.NewInjector()
	cfg, err := injector.NewConfigFromEnvironment()

	if err != nil {
		log.Printf("injector: error reading config from environment: %s\n", err)
		return
	}

	if err := i.Init(cfg, kubeClient); err != nil {
		log.Printf("injector: error initializing injector from config: %s", err)
		return
	}

	go func() {
		sig := <-sigchannel
		log.Printf("injector: received %s signal, shutting down\n", sig)
		cancel()
		sig = <-sigchannel
		log.Fatalf("injector: received %s signal during shutdown, exit immediately\n", sig)
	}()

	if err := i.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
