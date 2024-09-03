package main

import (
    "bytes"
    "context"
    "flag"
    "fmt"
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

// marshalRequest marshals a protobuf request object into a byte array.
func marshalRequest(helloReq *pb.HelloRequest) ([]byte, error) {
    binReq, err := proto.Marshal(helloReq)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %v", err)
    }
    return binReq, nil
}

// sendRequest sends an HTTP POST request with the given byte array and returns the response body as a byte array.
func sendRequest(addr string, binReq []byte, timeout int) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "POST", addr, bytes.NewBuffer(binReq))
    if err != nil {
        return nil, fmt.Errorf("failed to create HTTP request: %v", err)
    }
    req.Header.Set("Content-Type", "application/octet-stream")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send HTTP request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
    }

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    return respBody, nil
}

// unmarshalResponse unmarshalls a byte array into a protobuf response object.
func unmarshalResponse(respBody []byte) (*pb.HelloResponse, error) {
    var helloResp pb.HelloResponse
    if err := proto.Unmarshal(respBody, &helloResp); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %v", err)
    }
    return &helloResp, nil
}

func main() {
    flag.Parse()

    helloReq := &pb.HelloRequest{
        Name: *name,
    }

    binReq, err := marshalRequest(helloReq)
    if err != nil {
        log.Fatalf("Error marshaling request: %v", err)
    }

    respBody, err := sendRequest(*addr, binReq, *timeout)
    if err != nil {
        log.Fatalf("Error sending request: %v", err)
    }

    helloResp, err := unmarshalResponse(respBody)
    if err != nil {
        log.Fatalf("Error unmarshaling response: %v", err)
    }

    log.Printf("Greeting: %s", helloResp.GetText())
    fmt.Println(helloResp.GetText())
}
