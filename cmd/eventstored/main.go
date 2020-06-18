package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreasM009/eventstore/pkg/eventstored/runtime"
)

func main() {
	r := runtime.NewRuntime()

	if err := r.FromFlags(); err != nil {
		fmt.Println(err)
	}

	if err := r.Start(); err != nil {
		log.Println(err)
	}

	// Block, until SIGINT (Ctrl+C)
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Block until we receive our signal.
	<-c

	fmt.Println("User cancelled execution")
	os.Exit(0)
}
