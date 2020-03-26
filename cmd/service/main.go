package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/AndreasM009/eventstore-go/store"
	"github.com/AndreasM009/eventstore-service-go/pkg/eventstore"

	"github.com/AndreasM009/eventstore-service-go/pkg/config/standalone"
	"github.com/AndreasM009/eventstore-service-go/pkg/http"
)

func main() {
	// Flags
	portFlag := flag.Int("port", 5000, "Server port to use")
	configFilePathFlag := flag.String("config", "", "path to config file")

	flag.Parse()

	// Configuration
	cp := standalone.NewStandalone(*configFilePathFlag)

	configuration, err := cp.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Metadata for EventStore
	metadata := store.Metadata{
		Properties: map[string]string{},
	}

	for _, v := range configuration.Spec.Metadata {
		metadata.Properties[v.Name] = v.Value
	}

	// Registry
	registry := eventstore.NewRegistry()

	// create and init eventstore
	storage, err := registry.Create(configuration.Spec.Type)
	if err != nil {
		log.Fatal(err)
		return
	}

	if err := storage.Init(metadata); err != nil {
		log.Fatal(err)
		return
	}

	// Server
	s := http.NewServer(*portFlag, storage)
	s.StartNonBlocking()

	fmt.Printf("Server started on port: %d", *portFlag)

	// Block, until SIGINT (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	fmt.Println("User cancelled execution")
	os.Exit(0)
}
