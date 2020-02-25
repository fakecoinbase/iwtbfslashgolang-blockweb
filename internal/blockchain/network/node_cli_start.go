package network

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

func (nodeCLI *NodeCLI) startNode(nodeID, minerAddress string) {
	if len(minerAddress) > 0 {
		// TODO: Validate address
		//if ValidateAddress(minerAddress) {
		//	fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		//} else {
		//	log.Panic("Wrong miner address!")
		//}
	}

	bootNode(nodeID, minerAddress)
}
