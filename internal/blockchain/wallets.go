package blockchain

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (wallets *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := string(wallet.GetAddress()[:])

	wallets.Wallets[address] = wallet

	return address
}

func (wallets *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (wallets Wallets) GetWallet(address string) Wallet {
	wallet := wallets.Wallets[address]

	if wallet == nil {
		log.Panic("Address not found in your wallet")
	}

	return *wallet
}

func NewWallets(nodeID string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	// TODO: Error handling
	wallets.loadFromFile(nodeID)

	return &wallets, nil
}

func (wallets *Wallets) loadFromFile(nodeID string) error {
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	// TODO: Error handling
	fileContent, _ := ioutil.ReadFile(walletFile)

	var loadedWallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))

	// TODO: Error handling
	decoder.Decode(&loadedWallets)

	wallets.Wallets = loadedWallets.Wallets

	return nil
}

func (wallets Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(walletFile, nodeID)

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	// TODO: Error handling
	encoder.Encode(wallets)

	// TODO: Error handling
	ioutil.WriteFile(walletFile, content.Bytes(), 0644)
}
