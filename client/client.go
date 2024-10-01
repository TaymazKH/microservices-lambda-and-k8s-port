package main

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/proto"

    pb "main/genproto"
)

var (
    addr    *string
    timeout *int
)

const (
    greeterService = "greeter"
    sayHelloRPC    = "say-hello"
    sayByeRPC      = "say-bye"
)

// determineMessageType chooses the correct message type to initialize.
func determineMessageType(rpcName string) proto.Message {
    var msg proto.Message
    switch rpcName {
    case sayHelloRPC:
        msg = &pb.HelloResponse{}
    default:
        msg = &pb.ByeResponse{}
    }
    return msg
}

// SayHello represents the Greeter/SayHello RPC.
// custom headers can be sent and received.
func SayHello(helloRequest *pb.HelloRequest, header *http.Header) (*pb.HelloResponse, error) {
    binReq, err := marshalRequest(helloRequest)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*addr, greeterService, sayHelloRPC, &binReq, header, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, sayHelloRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.HelloResponse), nil
}

// SayBye represents the Greeter/SayBye RPC.
// custom headers can be sent and received.
func SayBye(byeRequest *pb.ByeRequest, header *http.Header) (*pb.ByeResponse, error) {
    binReq, err := marshalRequest(byeRequest)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*addr, greeterService, sayByeRPC, &binReq, header, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, sayByeRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ByeResponse), nil
}

// init loads the addr and timeout variables.
func init() {
    a, ok := os.LookupEnv("GREETER_SERVICE_ADDR")
    if !ok {
        log.Fatal("GREETER_SERVICE_ADDR environment variable not set")
    }
    addr = &a

    t, ok := os.LookupEnv("GREETER_SERVICE_TIMEOUT")
    if !ok {
        t := 5
        timeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil {
            t = 5
            timeout = &t
        } else {
            timeout = &t
        }
    }
}

// sendRequest sends an HTTP POST request with the given byte array and returns the response body as a byte array.
func sendRequest(addr, serviceName, rpcName string, binReq *[]byte, headers *http.Header, timeout int) ([]byte, *http.Header, error) {
    req, err := http.NewRequestWithContext(context.Background(), "POST", addr+"/"+serviceName, bytes.NewBuffer(*binReq))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
    }

    if headers != nil {
        req.Header = *headers
    }
    req.Header.Set("rpc-name", rpcName)
    req.Header.Set("content-type", "application/x-protobuf")

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
