package farmer

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

import (
	"github.com/alecthomas/kong"
)

var farmerCli struct {
	Start startFarmerCmd `cmd optional help:"Start a farmer server."`
}

func Run() {
	context := kong.Parse(&farmerCli)

	err := context.Run()
	context.FatalIfErrorf(err)
}
