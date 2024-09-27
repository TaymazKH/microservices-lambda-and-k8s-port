package main

import (
    "log"

    pb "main/genproto"
)

// handleSayHello processes the HelloRequest and returns a HelloResponse
func handleSayHello(helloRequest *pb.HelloRequest) (*pb.HelloResponse, error) {
    log.Printf("Received: %v", helloRequest.GetName())

    helloResp := &pb.HelloResponse{
        Text: "Hello " + helloRequest.GetName(),
    }

    return helloResp, nil
}

// handleSayBye processes the ByeRequest and returns a ByeResponse
func handleSayBye(byeRequest *pb.ByeRequest) (*pb.ByeResponse, error) {
    log.Printf("Received: %v", byeRequest.GetName())

    byeResp := &pb.ByeResponse{
        Text: "Bye " + byeRequest.GetName(),
    }

    return byeResp, nil
}
