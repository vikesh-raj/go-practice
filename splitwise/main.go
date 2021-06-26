package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vikesh-raj/go-practice/splitwise/server"
)

func run() error {
	config := flag.String("config", "config.yaml", "config file for the application server")
	port := flag.Int("port", 8192, "port on which the server should run")
	flag.Parse()

	fmt.Println("Starting server on port", *port)

	opts, err := server.GetOptsFromConfig(*config)
	if err != nil {
		return fmt.Errorf("Unable to read config : %v", err)
	}
	a, err := server.CreateApplication(opts)
	if err != nil {
		fmt.Printf("Unable to create application : %v\n", err)
	}

	err = a.StartServer(*port)
	if err != nil {
		return fmt.Errorf("Unable to start server : %v", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
