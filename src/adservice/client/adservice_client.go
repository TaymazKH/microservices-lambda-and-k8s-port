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

func main() {
    flag.Parse()

    adRequest := &pb.AdRequest{
        ContextKeys: []string{*contextKey},
    }

    binReq, err := proto.Marshal(adRequest)
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

    var adResponse pb.AdResponse
    if err := proto.Unmarshal(respBody, &adResponse); err != nil {
        log.Fatalf("Failed to decode response: %v", err)
    }

    for _, ad := range adResponse.GetAds() {
        log.Printf("Ad: %s", ad.GetText())
    }

    fmt.Println("Ad retrieval complete.")
}
