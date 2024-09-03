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
    addr       = flag.String("addr", "http://localhost:8080", "The address to connect to")
    contextKey = flag.String("contextKey", "camera", "Context key for ad request")
    timeout    = flag.Int("timeout", 5, "Timeout in seconds")
)

func marshalRequest(adRequest *pb.AdRequest) ([]byte, error) {
    return proto.Marshal(adRequest)
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

func unmarshalResponse(respBody []byte) (*pb.AdResponse, error) {
    var adResponse pb.AdResponse
    if err := proto.Unmarshal(respBody, &adResponse); err != nil {
        return nil, err
    }
    return &adResponse, nil
}

func main() {
    flag.Parse()

    adRequest := &pb.AdRequest{
        ContextKeys: []string{*contextKey},
    }

    binReq, err := marshalRequest(adRequest)
    if err != nil {
        log.Fatalf("Error marshaling request: %v", err)
    }

    respBody, err := sendRequest(*addr, binReq, *timeout)
    if err != nil {
        log.Fatalf("Error sending request: %v", err)
    }

    adResponse, err := unmarshalResponse(respBody)
    if err != nil {
        log.Fatalf("Error unmarshaling response: %v", err)
    }

    for _, ad := range adResponse.GetAds() {
        log.Printf("Ad: %s", ad.GetText())
    }

    fmt.Println("Ad retrieval complete.")
}
