package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwilliams-un/mud/player"
)

func shutdown(shutdownSignal os.Signal) {
	fmt.Printf("Shutdown (%s received)\n", shutdownSignal)
}

func main() {
	fmt.Println("MUD Start")

	// Capture SIGINT and SIGTERM signals and call shutdown to terminate cleanly
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		shutdown(sig)
		os.Exit(0)
	}()

	ln, err := net.Listen("tcp", ":4000")
	if err != nil {
		fmt.Println("Failed to listen for connections", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("Listening on %s\n", ln.Addr())

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Error handling new connection", err)
			} else {
				player.ConnectPlayer(conn)
			}
		}
	}()

	// Main loop
	for {

	}
}
