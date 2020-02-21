package dns_seed

/*
 * Copyright 2020 Information Wants To Be Free
 * Visit: https://github.com/iwtbf
 *
 * This project is licensed under the terms of the Apache 2.0 License.
 */

type DNSSeed struct {
}

func (dnsSeed *DNSSeed) Run() {
	// TODO: Error handling
	//lis,_ := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	//grpcServer := grpc.NewServer()
	//pb.RegisterRouteGuideServer(grpcServer, &routeGuideServer{})
	//grpcServer.Serve(lis)
}

func NewDNSSeed() *DNSSeed {
	return &DNSSeed{}
}
