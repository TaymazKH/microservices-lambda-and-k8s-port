package main

import (
    "log"

    pb "main/genproto"
)

// sayHello processes the HelloRequest and returns a HelloResponse
func sayHello(helloRequest *pb.HelloRequest, headers *map[string]string) (*pb.HelloResponse, error) {
    log.Printf("Received: %v", helloRequest.GetName())

    helloResp := &pb.HelloResponse{
        Text: "Hello " + helloRequest.GetName(),
    }

    return helloResp, nil
}

// sayBye processes the ByeRequest and returns a ByeResponse
func sayBye(byeRequest *pb.ByeRequest, headers *map[string]string) (*pb.ByeResponse, error) {
    log.Printf("Received: %v", byeRequest.GetName())

    byeResp := &pb.ByeResponse{
        Text: "Bye " + byeRequest.GetName(),
    }

    return byeResp, nil
}
