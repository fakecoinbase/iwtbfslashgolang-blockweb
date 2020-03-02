package dns_seed

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	knownEntryServers = []string{"/ip4/127.0.0.1/tcp/2000"}
	logger            = log.Logger("dns-seed")
)

type dnsSeed struct {
	UnimplementedFarmerServer
}

func (farmer *dnsSeed) RequestSeed(ctx context.Context, seedRequest *SeedRequest) (*SeedReply, error) {
	randomSeed := knownEntryServers[rand.Intn(len(knownEntryServers))]

	logger.Debugf("RequestSeed -> SeedReply{Seed: %s}", randomSeed)

	return &SeedReply{Seed: randomSeed}, nil
}

func startListening(listener net.Listener, grpcServer *grpc.Server) {
	logger.Info("Server startup success..")

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func bootDNSSeed(port int16) {
	logger.Infof("Booting DNS seed on port %d..", port)

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	RegisterFarmerServer(grpcServer, &dnsSeed{})

	go startListening(listener, grpcServer)

	logger.Info("Done.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	logger.Info("Received signal, shutting down...")
}
