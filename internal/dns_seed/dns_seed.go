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
	"google.golang.org/grpc"
	"math/rand"
	"net"
)

var (
	knownEntryServers = []string{"127.0.0.1:3000"}
)

type dnsSeed struct {
	UnimplementedFarmerServer
}

func (farmer *dnsSeed) RequestSeed(ctx context.Context, seedRequest *SeedRequest) (*SeedReply, error) {
	return &SeedReply{Seed: knownEntryServers[rand.Intn(len(knownEntryServers))]}, nil
}

func bootDNSSeed(port int) {
	fmt.Printf("Starting DNS seed on port %d\n", port)

	// TODO: Error handling
	listener, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	grpcServer := grpc.NewServer()
	RegisterFarmerServer(grpcServer, &dnsSeed{})

	// TODO: Background, it is blocking
	grpcServer.Serve(listener)

	fmt.Printf("Starting DNS seed on port %d\n", port)
}
