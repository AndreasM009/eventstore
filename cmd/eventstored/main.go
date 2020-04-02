package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/AndreasM009/eventstore-service-go/pkg/runtime"
)

func main() {
	r := runtime.NewRuntime()

	if err := r.FromFlags(); err != nil {
		fmt.Println(err)
		return
	}

	if err := r.Start(); err != nil {
		return
	}

	// Block, until SIGINT (Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	fmt.Println("User cancelled execution")
	os.Exit(0)
}
