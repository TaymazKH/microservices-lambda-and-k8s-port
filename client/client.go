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
    addr = flag.String("addr", "localhost:50051", "the address to connect to")
    name = flag.String("name", "world", "Name to greet")
)

func main() {
    flag.Parse()

    // Create the HelloRequest object
    req := &pb.HelloRequest{Name: *name}

    // Encode the HelloRequest object to binary (Protocol Buffers format)
    binReq, err := proto.Marshal(req)
    if err != nil {
        log.Fatalf("Failed to encode request: %v", err)
    }

    // Create the HTTP request
    httpReq, err := http.NewRequestWithContext(context.Background(), "POST", *addr, bytes.NewBuffer(binReq))
    if err != nil {
        log.Fatalf("Failed to create HTTP request: %v", err)
    }
    httpReq.Header.Set("Content-Type", "application/octet-stream")

    // Send the HTTP request
    client := &http.Client{Timeout: time.Second}
    resp, err := client.Do(httpReq)
    if err != nil {
        log.Fatalf("Failed to send HTTP request: %v", err)
    }
    defer resp.Body.Close()

    // Check if the response status is OK
    if resp.StatusCode != http.StatusOK {
        log.Fatalf("Received non-OK response: %s", resp.Status)
    }

    // Read the entire response body into a byte slice
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response body: %v", err)
    }

    // Decode the HelloResponse object from the binary response
    var helloResp pb.HelloResponse
    if err := proto.Unmarshal(respBody, &helloResp); err != nil {
        log.Fatalf("Failed to decode response: %v", err)
    }

    // Print the greeting from the response
    log.Printf("Greeting: %s", helloResp.GetText())
}
