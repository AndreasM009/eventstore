package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	eventstoreclient "github.com/AndreasM009/eventstore-service-go/pkg/client/clientset/versioned"
	"github.com/AndreasM009/eventstore-service-go/pkg/factories"
	"github.com/AndreasM009/eventstore-service-go/pkg/operator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigchannel := make(chan os.Signal)
	signal.Notify(sigchannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	config := factories.CreateKubeConfig()
	kubeClient := factories.CreateKubeClient()
	eventStoreClient, err := eventstoreclient.NewForConfig(config)

	if err != nil {
		log.Printf("Failed to create Eventstore client %v\n", err)
		return
	}

	operator := operator.NewOperator(eventStoreClient, kubeClient)

	done, err := operator.Run(ctx)
	if err != nil {
		log.Println("Failed to start operator")
		return
	}

	select {
	case <-sigchannel:
		cancel()
		log.Println("Operator received stop signal")
		<-time.After(2 * time.Second)
	case <-done:
		cancel()
	}
}
