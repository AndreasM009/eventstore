package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	eventstoreclient "github.com/AndreasM009/eventstore-service-go/pkg/client/clientset/versioned"
	"github.com/AndreasM009/eventstore-service-go/pkg/factories"
	"github.com/AndreasM009/eventstore-service-go/pkg/operator"
	"github.com/AndreasM009/eventstore-service-go/pkg/operator/http"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

func main() {
	apiPortFlags := flag.Int("port", 5000, "api server's port")
	flag.Parse()

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

	extensionClient, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create ApiExtension client %v\n", err)
		return
	}

	operator := operator.NewOperator(eventStoreClient, kubeClient, extensionClient)

	err = operator.InitCustomResourceDefinitions()
	if err != nil {
		log.Printf("Error creating CustomResoiurceDefinition: %s\n", err)
		return
	}

	done, err := operator.Run(ctx)
	if err != nil {
		log.Println("Failed to start operator")
		return
	}

	server := http.NewServer(*apiPortFlags, eventStoreClient)

	server.StartNonBlocking()

	select {
	case <-sigchannel:
		cancel()
		log.Println("Operator received stop signal")
		<-time.After(2 * time.Second)
	case <-done:
		cancel()
	}
}
