package client

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	"github.com/ipfs/go-log"
	"github.com/iwtbf/golang-blockweb/internal/farmer"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

var (
	clientLogger         = log.Logger("farmer-client")
	knownFarmerAddresses = []string{
		"localhost:10000",
	}
)

func selectRandomFarmerAddress(ignoredFarmerAddresses []string) string {
	farmerAddresses := make([]string, len(knownFarmerAddresses))
	copy(farmerAddresses, knownFarmerAddresses)

	for outer := 0; outer < len(ignoredFarmerAddresses); outer++ {
		for iterator := 0; iterator < len(farmerAddresses); iterator++ {
			if ignoredFarmerAddresses[outer] == farmerAddresses[iterator] {
				farmerAddresses[len(farmerAddresses)-1], farmerAddresses[iterator] = farmerAddresses[iterator], farmerAddresses[len(farmerAddresses)-1]
				farmerAddresses = farmerAddresses[:len(farmerAddresses)-1]
			}
		}
	}

	return farmerAddresses[rand.Intn(len(farmerAddresses))]
}

func requestCoreRelayFromFarmer(address, hostAddress string) (string, []byte) {
	clientLogger.Debugf("Trying random farmer: %s", address)

	dialCtx, dialCancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer dialCancelFunc()
	clientConn, err := grpc.DialContext(dialCtx, address, grpc.WithBlock())
	if err != nil {
		clientLogger.Warning(err)
		return "", nil
	}

	dialCancelFunc()
	defer clientConn.Close()

	client := farmer.NewFarmerClient(clientConn)

	requestCtx, requestCancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer requestCancelFunc()

	seed, err := client.RequestSeed(requestCtx, &farmer.SeedRequest{HostAddress: hostAddress})
	if err != nil {
		clientLogger.Warning(err)
		return "", nil
	}

	requestCancelFunc()

	return seed.Address, seed.PublicKeyHash
}

func RequestRandomCoreRelayInformation(hostAddress string) (string, []byte) {
	clientLogger.Info("Looking out for random core relay.")

	var testedFarmerAddresses []string

	var relayAddress string
	var relayPublicKeyHash []byte

	for i := 0; relayAddress == "" && i < 3; i++ {
		farmerAddress := selectRandomFarmerAddress(testedFarmerAddresses)
		relayAddress, relayPublicKeyHash = requestCoreRelayFromFarmer(farmerAddress, hostAddress)

		if relayAddress == "" {
			testedFarmerAddresses = append(testedFarmerAddresses, farmerAddress)
		}
	}

	clientLogger.Infof("Found core relay with address: %s", relayAddress)

	return relayAddress, relayPublicKeyHash
}
