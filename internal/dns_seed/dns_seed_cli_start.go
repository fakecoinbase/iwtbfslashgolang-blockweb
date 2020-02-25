package dns_seed

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import "github.com/ipfs/go-log"

func (cmd *startDnsSeedCmd) Run() error {
	if cmd.Level != "" {
		level, err := log.LevelFromString(cmd.Level)
		if err != nil {
			panic(err)
		}

		log.SetAllLoggers(level)
	}

	bootDNSSeed(cmd.Port)

	return nil
}