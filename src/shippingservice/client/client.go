package client

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
    addr    = flag.String("addr", "localhost:50051", "The address to connect to")
    timeout = flag.Int("timeout", 5, "Timeout in seconds")
)

const (
    shippingService = "shipping-service"
    getQuoteRPC     = "get-quote"
    shipOrderRPC    = "ship-order"
)

func GetQuote(getQuoteRequest *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
    binReq, err := marshalRequest(getQuoteRequest)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/%s/%s", shippingService, getQuoteRPC)
    respBody, header, err := sendRequest(*addr, path, binReq, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, path)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.GetQuoteResponse), nil
}

func ShipOrder(shipOrderRequest *pb.ShipOrderRequest) (*pb.ShipOrderResponse, error) {
    binReq, err := marshalRequest(shipOrderRequest)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/%s/%s", shippingService, shipOrderRPC)
    respBody, header, err := sendRequest(*addr, path, binReq, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, path)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ShipOrderResponse), nil
}

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

func marshalRequest(msg proto.Message) ([]byte, error) {
    binReq, err := proto.Marshal(msg)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    return binReq, nil
}

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
        case fmt.Sprintf("/%s/%s", shippingService, getQuoteRPC):
            msg = &pb.GetQuoteResponse{}
        default:
            msg = &pb.ShipOrderResponse{}
        }

        if err := proto.Unmarshal(respBody, msg); err != nil {
            return nil, fmt.Errorf("failed to unmarshal response: %w", err)
        }
        return &msg, nil
    } else {
        return nil, status.Error(grpcCode, string(respBody))
    }
}
