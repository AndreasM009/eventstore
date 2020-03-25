package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/AndreasM009/eventstore-go/store/azure/tablestorage"
	"github.com/AndreasM009/eventstore-service-go/pkg/http"
)

func main() {
	// Flags
	portFlag := flag.Int("port", 5000, "Server port to use")
	storageAccountNameFlag := flag.String("storageaccountname", "", "name of azure storage account")
	storageAccountKeyFlag := flag.String("storageaccountkey", "", "key of azure storage account")

	flag.Parse()

	if *storageAccountNameFlag == "" || *storageAccountKeyFlag == "" {
		flag.Usage()
		return
	}

	// EventStore
	eventstore := tablestorage.NewStore(*storageAccountNameFlag, *storageAccountKeyFlag, "")
	if err := eventstore.Init(); err != nil {
		log.Fatal(err)
		return
	}

	// Server
	s := http.NewServer(*portFlag, eventstore)
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
