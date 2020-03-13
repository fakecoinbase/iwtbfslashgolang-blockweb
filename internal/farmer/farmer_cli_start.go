package farmer

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/ipfs/go-log"
	"github.com/whyrusleeping/go-logging"
)

type startFarmerCmd struct {
	CertFile string `flag type:"existingfile" help:"Path to the cert file."`
	KeyFile  string `flag type:"existingfile" help:"Path to the ECDSA private key file."`
	Port     int16  `flag optional help:"The servers listening Port." default:"10000"`
	Level    string `flag optional help:"One of github.com/whyrusleeping/go-logging#LogLevel." default:"INFO"`
}

func (startFarmer *startFarmerCmd) Run() error {
	if startFarmer.Level != "" {
		level, err := logging.LogLevel(startFarmer.Level)
		if err != nil {
			level = logging.INFO
		}

		log.SetAllLoggers(level)
	}

	bootFarmer(startFarmer.CertFile, startFarmer.KeyFile, startFarmer.Port)

	return nil
}
