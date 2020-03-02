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
	"github.com/whyrusleeping/go-logging"
	"os"
)

type startNodeCmd struct {
	KeyFile string `arg type:"existingfile" help:"Path to the ECDSA private key file."`
	Port    int16  `optional help:"The servers listening Port." default:"2000"`
	Level   string `optional help:"One of github.com/whyrusleeping/go-logging#LogLevel." default:"INFO"`
	Lazy    bool   `optional help:"A 'proxy' node which does neither verify the blockchain nor act as a relay."`
}

func loadKeyPairFromCertFile(certFile string) keygen.KeyPair {
	keyPairs := keygen.LoadKeyPairsFromExistingFile(certFile)

	if len(keyPairs) == 0 {
		logger.Error("No ECDSA private keys found in file!")
		os.Exit(1)
	}

	if len(keyPairs) > 1 {
		logger.Error("Multiple ECDSA private keys found in file!")
		os.Exit(1)
	}

	return keyPairs[0]
}

func (cmd *startNodeCmd) Run() error {
	if cmd.Level != "" {
		level, err := logging.LogLevel(cmd.Level)
		if err != nil {
			level = logging.INFO
		}

		log.SetAllLoggers(level)
	}

	keyPair := loadKeyPairFromCertFile(cmd.KeyFile)

	if cmd.Lazy {
		logger.Error("Lazy node not implement yet!")
		os.Exit(1)
	} else {
		bootRelay(cmd.Port, keyPair)
	}

	return nil
}
