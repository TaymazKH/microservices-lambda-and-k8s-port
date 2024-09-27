package main

import (
    "flag"
    "fmt"
    "log"

    pb "main/genproto"
)

var (
    name = flag.String("name", "world", "Name to greet")
)

func main() {
    flag.Parse()

    helloReq := &pb.HelloRequest{
        Name: *name,
    }

    helloResp, err := SayHello(helloReq)
    if err != nil {
        log.Fatalf("Error calling SayHello RPC: %v", err)
    }

    log.Printf("Greeting: %s", helloResp.GetText())

    byeReq := &pb.ByeRequest{
        Name: *name,
    }

    byeResp, err := SayBye(byeReq)
    if err != nil {
        log.Fatalf("Error calling SayBye RPC: %v", err)
    }

    log.Printf("Farewell: %s", byeResp.GetText())

    fmt.Println(helloResp.GetText(), "-", byeResp.GetText())
}
