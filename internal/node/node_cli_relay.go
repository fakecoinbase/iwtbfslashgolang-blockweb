package node

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/ipfs/go-log"
	"github.com/iwtbf/golang-blockweb/internal/keygen"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/whyrusleeping/go-logging"
	"google.golang.org/grpc/credentials"
	"os"
)

type relayNodeCmd struct {
	CertFile string `flag type:"existingfile" help:"Path to the cert file."`
	KeyFile  string `flag type:"existingfile" help:"Path to the ECDSA private key file."`
	Port     int16  `optional help:"The servers listening Port." default:"2000"`
	Level    string `optional help:"One of github.com/whyrusleeping/go-logging#LogLevel." default:"INFO"`
	Lazy     bool   `optional help:"A 'proxy' node which does neither verify the blockchain nor act as a node."`
}

func loadKeyPairFromCertFile(certFile string) crypto.PrivKey {
	keyPairs := keygen.LoadKeyPairsFromExistingFile(certFile)

	if len(keyPairs) == 0 {
		logger.Error("No ECDSA private keys found in file!")
		os.Exit(1)
	}

	if len(keyPairs) > 1 {
		logger.Error("Multiple ECDSA private keys found in file!")
		os.Exit(1)
	}

	privKey, _, err := crypto.ECDSAKeyPairFromKey(&keyPairs[0].PrivateKey)
	if err != nil {
		panic(err)
	}

	return privKey
}

func (cmd *relayNodeCmd) Run() error {
	if cmd.Level != "" {
		level, err := logging.LogLevel(cmd.Level)
		if err != nil {
			level = logging.INFO
		}

		log.SetAllLoggers(level)
	}

	transportCredentials, err := credentials.NewServerTLSFromFile(cmd.CertFile, cmd.KeyFile)
	if err != nil {
		panic(err)
	}

	privKey := loadKeyPairFromCertFile(cmd.KeyFile)

	if cmd.Lazy {
		logger.Error("Lazy node not implement yet!")
		os.Exit(1)
	} else {
		bootRelay(cmd.Port, transportCredentials, privKey)
	}

	return nil
}
