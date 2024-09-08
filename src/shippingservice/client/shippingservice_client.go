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
    pb "main/genproto"
)

var (
    addr    = flag.String("addr", "http://localhost:8080", "The address to connect to")
    timeout = flag.Int("timeout", 5, "Timeout in seconds")
)

func marshalRequest(request *pb.GetQuoteRequest) ([]byte, error) {
    return proto.Marshal(request)
}

func sendRequest(addr string, binReq []byte, timeout int) ([]byte, error) {
    req, err := http.NewRequestWithContext(context.Background(), "POST", addr, bytes.NewBuffer(binReq))
    if err != nil {
        return nil, fmt.Errorf("failed to create HTTP request: %v", err)
    }
    req.Header.Set("Content-Type", "application/octet-stream")

    client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
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

func unmarshalResponse(respBody []byte) (*pb.GetQuoteResponse, error) {
    var response pb.GetQuoteResponse
    if err := proto.Unmarshal(respBody, &response); err != nil {
        return nil, err
    }
    return &response, nil
}

func main() {
    flag.Parse()

    req := &pb.GetQuoteRequest{
        Address: &pb.Address{
            StreetAddress: "Muffin Man",
            City:          "London",
            State:         "",
            Country:       "England",
        },
        Items: []*pb.CartItem{
            {
                ProductId: "23",
                Quantity:  1,
            },
            {
                ProductId: "46",
                Quantity:  3,
            },
        },
    }

    binReq, err := marshalRequest(req)
    if err != nil {
        log.Fatalf("Error marshaling request: %v", err)
    }

    respBody, err := sendRequest(*addr, binReq, *timeout)
    if err != nil {
        log.Fatalf("Error sending request: %v", err)
    }

    resp, err := unmarshalResponse(respBody)
    if err != nil {
        log.Fatalf("Error unmarshaling response: %v", err)
    }

    log.Printf("Quote: %d.%d", resp.CostUsd.GetUnits(), resp.CostUsd.GetNanos())
    fmt.Println("Quote retrieval complete.")
}
