package main

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    greeterAddr    *string
    greeterTimeout *int
)

const (
    greeterService = "greeter"
    sayHelloRPC    = "say-hello"
    sayByeRPC      = "say-bye"
)

// SayHello represents the Greeter/SayHello RPC.
// context can be sent as custom headers.
func SayHello(request *pb.HelloRequest, header *http.Header) (*pb.HelloResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*greeterAddr, greeterService, sayHelloRPC, &binReq, header, *greeterTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, sayHelloRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.HelloResponse), nil
}

// SayBye represents the Greeter/SayBye RPC.
// context can be sent as custom headers.
func SayBye(request *pb.ByeRequest, header *http.Header) (*pb.ByeResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*greeterAddr, greeterService, sayByeRPC, &binReq, header, *greeterTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, sayByeRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ByeResponse), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("GREETER_SERVICE_ADDR")
    if !ok {
        log.Fatal("GREETER_SERVICE_ADDR environment variable not set")
    }
    greeterAddr = &a

    t, ok := os.LookupEnv("GREETER_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        greeterTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            greeterTimeout = &t
        } else {
            greeterTimeout = &t
        }
    }
}
