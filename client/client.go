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

// SayHello sends a HelloRequest to server and returns a HelloResponse
func SayHello(helloRequest *pb.HelloRequest) (*pb.HelloResponse, error) {
    binReq, err := marshalRequest(helloRequest)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/%s/%s", greeterService, sayHelloRPC)
    respBody, header, err := sendRequest(*addr, path, binReq, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, path)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.HelloResponse), nil
}

// SayBye sends a ByeRequest to server and returns a ByeResponse
func SayBye(byeRequest *pb.ByeRequest) (*pb.ByeResponse, error) {
    binReq, err := marshalRequest(byeRequest)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/%s/%s", greeterService, sayByeRPC)
    respBody, header, err := sendRequest(*addr, path, binReq, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, path)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ByeResponse), nil
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

// marshalRequest marshals a proto message object into a byte array.
func marshalRequest(msg proto.Message) ([]byte, error) {
    binReq, err := proto.Marshal(msg)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    return binReq, nil
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
        case fmt.Sprintf("/%s/%s", greeterService, sayHelloRPC):
            msg = &pb.HelloResponse{}
        default:
            msg = &pb.ByeResponse{}
        }

        if err := proto.Unmarshal(respBody, msg); err != nil {
            return nil, fmt.Errorf("failed to unmarshal response: %w", err)
        }
        return &msg, nil
    } else {
        return nil, status.Error(grpcCode, string(respBody))
    }
}
