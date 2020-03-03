package farmer

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func startListening(listener net.Listener, grpcServer *grpc.Server) {
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

// TODO: Add TLS
func bootFarmer(port int16) {
	farmerLogger.Infof("Booting farmer on port %d..", port)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	RegisterFarmerServer(grpcServer, &farmer{})

	go startListening(listener, grpcServer)

	farmerLogger.Info("Done.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	farmerLogger.Info("Received signal, shutting down...")
}
