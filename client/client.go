package main

import (
    "bytes"
    "context"
    "flag"
    "io"
    "log"
    "net/http"
    "time"

    "google.golang.org/protobuf/proto"
    pb "main/hello"
)

var (
    addr    = flag.String("addr", "localhost:50051", "the address to connect to")
    name    = flag.String("name", "world", "Name to greet")
    timeout = flag.Int("timeout", 5, "Timeout in seconds")
)

func main() {
    flag.Parse()

    helloReq := &pb.HelloRequest{
        Name: *name,
    }

    binReq, err := proto.Marshal(helloReq)
    if err != nil {
        log.Fatalf("Failed to encode request: %v", err)
    }

    req, err := http.NewRequestWithContext(context.Background(), "POST", *addr, bytes.NewBuffer(binReq))
    if err != nil {
        log.Fatalf("Failed to create HTTP request: %v", err)
    }
    req.Header.Set("Content-Type", "application/octet-stream")

    client := &http.Client{Timeout: time.Duration(*timeout) * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Failed to send HTTP request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Fatalf("Received non-OK response: %s", resp.Status)
    }

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response body: %v", err)
    }

    var helloResp pb.HelloResponse
    if err := proto.Unmarshal(respBody, &helloResp); err != nil {
        log.Fatalf("Failed to decode response: %v", err)
    }

    log.Printf("Greeting: %s", helloResp.GetText())
}
