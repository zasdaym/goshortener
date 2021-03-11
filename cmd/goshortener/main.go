package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/zasdaym/goshortener/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signal := <-ch
		log.Printf("got signal: %v", signal)
		cancel()
	}()

	port := flag.String("port", "8080", "Port to listen")
	addr := ":" + *port
	srv := http.NewServer(addr)
	return srv.Serve(ctx)
}
