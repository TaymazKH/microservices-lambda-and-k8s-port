package main

import (
    "bytes"
    "context"
    "flag"
    "fmt"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "io"
    "net/http"
    "strconv"
    "time"

    "google.golang.org/protobuf/proto"
    pb "main/genproto"
)

var (
    addr       = flag.String("addr", "http://localhost:8080", "The address to connect to")
    contextKey = flag.String("contextKey", "camera", "Context key for ad request")
    timeout    = flag.Int("timeout", 5, "Timeout in seconds")
)

const (
    adService = "ad-service"
    getAdsRPC = "get-ads"
)

func GetAds(adRequest *pb.AdRequest) (*pb.AdResponse, error) {
    binReq, err := marshalRequest(adRequest)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/%s/%s", adService, getAdsRPC)
    respBody, header, err := sendRequest(*addr, path, binReq, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, path)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.AdResponse), nil
}

// marshalRequest marshals a proto message object into a byte array.
func marshalRequest(msg proto.Message) ([]byte, error) {
    binReq, err := proto.Marshal(msg)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    return binReq, nil
}

// sendRequest sends an HTTP POST request with the given byte array and returns the response body as a byte array.
func sendRequest(addr, path string, binReq []byte, timeout int) ([]byte, *http.Header, error) {
    req, err := http.NewRequestWithContext(context.Background(), "POST", addr+path, bytes.NewBuffer(binReq))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
    }
    req.Header.Set("Content-Type", "application/octet-stream")

    client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to send HTTP request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, nil, fmt.Errorf("received non-OK response: %s", resp.Status)
    }

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, nil, fmt.Errorf("failed to read response body: %w", err)
    }

    return respBody, &resp.Header, nil
}

// unmarshalResponse unmarshalls a byte array into a proto message object.
func unmarshalResponse(respBody []byte, header *http.Header, path string) (*proto.Message, error) {
    if header.Get("Grpc-Code") == "" {
        return nil, fmt.Errorf("missing Grpc-Code header")
    }

    grpcCode, err := strconv.Atoi(header.Get("Grpc-Code"))
    if err != nil {
        return nil, fmt.Errorf("failed to parse Grpc-Code header: %w", err)
    }

    if grpcCode := codes.Code(grpcCode); grpcCode == codes.OK {
        var msg proto.Message
        switch path {
        default:
            msg = &pb.AdResponse{}
        }

        if err := proto.Unmarshal(respBody, msg); err != nil {
            return nil, fmt.Errorf("failed to unmarshal response: %w", err)
        }
        return &msg, nil
    } else {
        return nil, status.Error(grpcCode, string(respBody))
    }
}
