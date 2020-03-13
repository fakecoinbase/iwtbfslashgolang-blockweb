package farmer

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
	"google.golang.org/grpc/credentials"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	farmerLogger = log.Logger("farmer")
)

type farmer struct {
	UnimplementedFarmerServer
}

func (farmer *farmer) RequestSeed(ctx context.Context, seedRequest *SeedRequest) (*SeedReply, error) {
	coreRelay := knownCoreRelays[0]

	if len(knownCoreRelays) > 1 {
		coreRelays := prepareCoreRelays(seedRequest.HostAddress)
		coreRelay = coreRelays[rand.Intn(len(coreRelays))]
	}

	farmerLogger.Debugf("RequestSeed -> SeedReply{Address: %s, PublicKeyHash: %s}", coreRelay.address, coreRelay.publicKeyHash)

	return &SeedReply{Address: coreRelay.address, PublicKeyHash: []byte(coreRelay.publicKeyHash)}, nil
}

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

	// TODO: credentials.NewServerTLSFromFile()
	credentials.NewServerTLSFromFile()
	grpcServer := grpc.NewServer()
	RegisterFarmerServer(grpcServer, &farmer{})

	go startListening(listener, grpcServer)

	farmerLogger.Info("Done.")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	farmerLogger.Info("Received signal, shutting down...")
}
