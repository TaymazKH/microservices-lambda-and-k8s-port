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
    addr    = flag.String("addr", "localhost:50051", "the address to connect to")
    name    = flag.String("name", "world", "Name to greet")
    timeout = flag.Int("timeout", 5, "Timeout in seconds")
)

const (
    greeterService = "greeter"
    sayHelloRPC    = "say-hello"
    sayByeRPC      = "say-bye"
)

// sayHello sends a HelloRequest to server and returns a HelloResponse
func sayHello(helloRequest *pb.HelloRequest) (*pb.HelloResponse, error) {
    binReq, err := marshalRequest(helloRequest)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/%s/%s", greeterService, sayHelloRPC)
    respBody, err := sendRequest(*addr, path, binReq, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, path)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.HelloResponse), nil
}

// sayBye sends a ByeRequest to server and returns a ByeResponse
func sayBye(byeRequest *pb.ByeRequest) (*pb.ByeResponse, error) {
    binReq, err := marshalRequest(byeRequest)
    if err != nil {
        return nil, err
    }

    path := fmt.Sprintf("/%s/%s", greeterService, sayByeRPC)
    respBody, err := sendRequest(*addr, path, binReq, *timeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, path)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ByeResponse), nil
}

// sendRequest sends an HTTP POST request with the given byte array and returns the response body as a byte array.
func sendRequest(addr, path string, binReq []byte, timeout int) ([]byte, error) {
    req, err := http.NewRequestWithContext(context.Background(), "POST", addr+path, bytes.NewBuffer(binReq))
    if err != nil {
        return nil, fmt.Errorf("failed to create HTTP request: %w", err)
    }
    req.Header.Set("Content-Type", "application/octet-stream")

    client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send HTTP request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
    }

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    return respBody, nil
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
func unmarshalResponse(respBody []byte, path string) (*proto.Message, error) {
    var msg proto.Message
    switch path {
    case fmt.Sprintf("/%s/%s", greeterService, sayHelloRPC):
        msg = &pb.HelloResponse{}
    case fmt.Sprintf("/%s/%s", greeterService, sayByeRPC):
        msg = &pb.ByeResponse{}
    default:
        return nil, fmt.Errorf("unknown path: %s", path)
    }

    if err := proto.Unmarshal(respBody, msg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    return &msg, nil
}

func main() {
    flag.Parse()

    helloReq := &pb.HelloRequest{
        Name: *name,
    }

    helloResp, err := sayHello(helloReq)
    if err != nil {
        log.Fatalf("Error calling SayHello RPC: %v", err)
    }

    log.Printf("Greeting: %s", helloResp.GetText())

    byeReq := &pb.ByeRequest{
        Name: *name,
    }

    byeResp, err := sayBye(byeReq)
    if err != nil {
        log.Fatalf("Error calling SayBye RPC: %v", err)
    }

    log.Printf("Farewell: %s", byeResp.GetText())

    fmt.Println(helloResp.GetText(), "-", byeResp.GetText())
}
