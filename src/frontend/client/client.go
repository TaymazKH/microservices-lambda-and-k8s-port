package client

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "time"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/proto"

    pb "main/genproto"
)

const (
    defaultTimeout = 20
)

// determineMessageType chooses the correct message type to initialize.
func determineMessageType(rpcName string) proto.Message {
    var msg proto.Message
    switch rpcName {
    case listProductsRPC:
        msg = &pb.ListProductsResponse{}
    case getProductRPC:
        msg = &pb.Product{}
    case searchProductsRPC:
        msg = &pb.SearchProductsResponse{}
    case getCartRPC:
        msg = &pb.Cart{}
    case getQuoteRPC:
        msg = &pb.GetQuoteResponse{}
    case shipOrderRPC:
        msg = &pb.ShipOrderResponse{}
    case convertRPC:
        msg = &pb.Money{}
    case listRecommendationsRPC:
        msg = &pb.ListRecommendationsResponse{}
    case getAdsRPC:
        msg = &pb.AdResponse{}
    case getSupportedCurrenciesRPC:
        msg = &pb.GetSupportedCurrenciesResponse{}
    case placeOrderRPC:
        msg = &pb.PlaceOrderResponse{}
    default:
        /*
           CartService/AddItem
           CartService/EmptyCart
        */
        msg = &pb.Empty{}
    }
    return msg
}

// sendRequest sends an HTTP POST request with the given byte array and returns the response body as a byte array.
func sendRequest(addr, serviceName, rpcName string, binReq *[]byte, headers *http.Header, timeout int) ([]byte, *http.Header, error) {
    req, err := http.NewRequest(http.MethodPost, addr+"/"+serviceName, bytes.NewBuffer(*binReq))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
    }

    if headers != nil {
        req.Header = *headers
    }
    req.Header.Set("rpc-name", rpcName)
    req.Header.Set("content-type", "application/octet-stream")

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

// marshalRequest marshals a protobuf message into a byte array.
func marshalRequest(msg proto.Message) ([]byte, error) {
    binReq, err := proto.Marshal(msg)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    return binReq, nil
}

// unmarshalResponse unmarshalls a byte array into a protobuf message.
func unmarshalResponse(respBody []byte, header *http.Header, rpcName string) (*proto.Message, error) {
    if header.Get("grpc-status") == "" {
        return nil, fmt.Errorf("missing grpc-status header")
    }

    grpcStatus, err := strconv.Atoi(header.Get("grpc-status"))
    if err != nil {
        return nil, fmt.Errorf("failed to parse grpc-status header: %w", err)
    }

    if grpcStatus := codes.Code(grpcStatus); grpcStatus == codes.OK {
        msg := determineMessageType(rpcName)

        if err := proto.Unmarshal(respBody, msg); err != nil {
            return nil, fmt.Errorf("failed to unmarshal response: %w", err)
        }
        return &msg, nil
    } else {
        return nil, status.Error(grpcStatus, string(respBody))
    }
}
