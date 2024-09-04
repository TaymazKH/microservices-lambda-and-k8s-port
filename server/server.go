package main

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "os"

    "google.golang.org/protobuf/proto"
    pb "main/genproto"
)

// RequestContext represents the nested context of the request
type RequestContext struct {
    HTTP struct {
        Method string `json:"method"`
        Path   string `json:"path"`
    } `json:"http"`
}

// RequestData represents the structure of the incoming JSON string
type RequestData struct {
    Body            string            `json:"body"`
    Headers         map[string]string `json:"headers"`
    RequestContext  RequestContext    `json:"requestContext"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// ResponseData represents the structure of the outgoing JSON string
type ResponseData struct {
    StatusCode      int               `json:"statusCode"`
    Headers         map[string]string `json:"headers"`
    Body            string            `json:"body"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// handleSayHello processes the HelloRequest and returns a HelloResponse
func handleSayHello(helloRequest *pb.HelloRequest) (*pb.HelloResponse, error) {
    log.Printf("Received: %v", helloRequest.GetName())

    helloResp := &pb.HelloResponse{
        Text: "Hello " + helloRequest.GetName(),
    }

    return helloResp, nil
}

// handleSayBye processes the ByeRequest and returns a ByeResponse
func handleSayBye(byeRequest *pb.ByeRequest) (*pb.ByeResponse, error) {
    log.Printf("Received: %v", byeRequest.GetName())

    byeResp := &pb.ByeResponse{
        Text: "Bye " + byeRequest.GetName(),
    }

    return byeResp, nil
}

// handleRequest chooses the correct handler function to call
func handleRequest(msg proto.Message, reqData *RequestData) (proto.Message, error) {
    switch reqData.RequestContext.HTTP.Path {
    case "/greeter/say-hello":
        return handleSayHello(msg.(*pb.HelloRequest))
    case "/greeter/say-bye":
        return handleSayBye(msg.(*pb.ByeRequest))
    default:
        return nil, fmt.Errorf("unknown path: %s", reqData.RequestContext.HTTP.Path)
    }
}

// decodeRequest decodes the incoming JSON request into a protobuf message
func decodeRequest(request string) (*proto.Message, *RequestData, error) {
    var reqData RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        return nil, nil, fmt.Errorf("failed to parse request JSON: %w", err)
    }

    var binReqBody []byte
    if reqData.IsBase64Encoded {
        var err error
        binReqBody, err = base64.StdEncoding.DecodeString(reqData.Body)
        if err != nil {
            return nil, nil, fmt.Errorf("failed to decode base64 body: %w", err)
        }
    } else {
        binReqBody = []byte(reqData.Body)
    }

    var msg proto.Message
    switch reqData.RequestContext.HTTP.Path {
    case "/greeter/say-hello":
        msg = &pb.HelloRequest{}
    case "/greeter/say-bye":
        msg = &pb.ByeRequest{}
    default:
        return nil, nil, fmt.Errorf("unknown path: %s", reqData.RequestContext.HTTP.Path)
    }

    if err := proto.Unmarshal(binReqBody, msg); err != nil {
        return nil, nil, fmt.Errorf("failed to unmarshal request body: %w", err)
    }

    return &msg, &reqData, nil
}

// encodeResponse encodes the protobuf response into the outgoing JSON response
func encodeResponse(msg proto.Message) (string, error) {
    binRespBody, err := proto.Marshal(msg)
    if err != nil {
        return "", fmt.Errorf("failed to marshal response: %w", err)
    }

    encodedRespBody := base64.StdEncoding.EncodeToString(binRespBody) // Base64 encoding is optional.

    respData := ResponseData{
        StatusCode:      200,
        Headers:         map[string]string{"Content-Type": "application/octet-stream"},
        Body:            encodedRespBody, // Use `binRespBody` if not encoded.
        IsBase64Encoded: true,
    }

    jsonResponse, err := json.Marshal(respData)
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON response: %w", err)
    }

    return string(jsonResponse), nil
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    request, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalf("Failed to read from stdin: %v", err)
    }
    request = request[:len(request)-1] // Trim any trailing newline characters

    msg, reqData, err := decodeRequest(request)
    if err != nil {
        log.Fatalf("Error decoding request: %v", err)
    }

    helloResp, err := handleRequest(*msg, reqData)
    if err != nil {
        log.Fatalf("Handler error: %v", err)
    }

    response, err := encodeResponse(helloResp)
    if err != nil {
        log.Fatalf("Error encoding response: %v", err)
    }

    fmt.Println(response)
}
