package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"context"
	farmerClient "github.com/iwtbf/golang-blockweb/internal/farmer/client"
	core "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

const (
	coreRelayRendezvous = "core-relay-rendezvous"
	relayRendezvous     = "relay-rendezvous"
)

type DistributedHashTable struct {
	*dht.IpfsDHT
}

func (distributedHashTable *DistributedHashTable) fetchHashTableFromCoreRelay(host core.Host, relayAddress string, relayPublicKeyHash []byte) {
	multiaddr, err := multiaddr.NewMultiaddr(relayAddress)
	if err != nil {
		panic(err)
	}

	peerInfo, err := peer.AddrInfoFromP2pAddr(multiaddr)
	if err != nil {
		panic(err)
	}

	if err := host.Connect(context.Background(), *peerInfo); err != nil {
		panic(err)
	} else {
		logger.Debugf("Connection to core relay '%v' established.", *peerInfo)
	}
}

func (distributedHashTable *DistributedHashTable) synchronize(host core.Host) {
	logger.Debug("Initializing distributed hash table..")

	peerInfo := peer.AddrInfo{ID: host.ID(), Addrs: host.Addrs()}
	addresses, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	primaryHostAddress := addresses[0].String()

	relayAddress, relayPublicKeyHash := farmerClient.RequestRandomCoreRelayInformation(primaryHostAddress)

	if relayAddress != primaryHostAddress {
		distributedHashTable.fetchHashTableFromCoreRelay(host, relayAddress, relayPublicKeyHash)
	} else {
		// TODO: Be the first to announce coreRelayRendezvous
		logger.Warning("I am the genesis relay - I won't bow to anyone!")
	}

	logger.Debug("Hash table initialized.")
}

func (distributedHashTable *DistributedHashTable) announce() {
	logger.Info("Announcing new relay..")
	routingDiscovery := discovery.NewRoutingDiscovery(distributedHashTable)
	discovery.Advertise(context.Background(), routingDiscovery, relayRendezvous)
	logger.Debug("Successfully announced.")
}

func newDistributedHashTable(relay *relay) *DistributedHashTable {
	kademliaDHT, err := dht.New(context.Background(), relay.host)
	if err != nil {
		panic(err)
	}

	return &DistributedHashTable{IpfsDHT: kademliaDHT}
}
