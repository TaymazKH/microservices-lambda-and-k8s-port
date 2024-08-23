package main

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "os"

    "google.golang.org/protobuf/proto"
    pb "main/hello"
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

// decodeRequest decodes the incoming JSON request into the protobuf message
func decodeRequest(request string) (*pb.HelloRequest, RequestData, error) {
    var reqData RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        log.Fatalf("Failed to parse request JSON: %v", err)
    }

    var binReqBody []byte
    if reqData.IsBase64Encoded {
        var err error
        binReqBody, err = base64.StdEncoding.DecodeString(reqData.Body)
        if err != nil {
            log.Fatalf("Failed to decode base64 body: %v", err)
        }
    } else {
        binReqBody = []byte(reqData.Body)
    }

    var helloReq pb.HelloRequest
    if err := proto.Unmarshal(binReqBody, &helloReq); err != nil {
        log.Fatalf("Failed to unmarshal request body: %v", err)
    }

    return &helloReq, reqData, nil
}

// sayHelloHandler processes the HelloRequest and returns a HelloResponse
func sayHelloHandler(helloReq *pb.HelloRequest) (*pb.HelloResponse, error) {
    log.Printf("Received: %v", helloReq.GetName())

    helloResp := &pb.HelloResponse{
        Text: "Hello " + helloReq.GetName(),
    }

    return helloResp, nil
}

// encodeResponse encodes the protobuf response into the outgoing JSON response
func encodeResponse(helloResp *pb.HelloResponse) (string, error) {
    binRespBody, err := proto.Marshal(helloResp)
    if err != nil {
        log.Fatalf("Failed to marshal response: %v", err)
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
        log.Fatalf("Failed to marshal JSON response: %v", err)
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

    helloReq, _, err := decodeRequest(request)
    if err != nil {
        log.Fatalf("Error decoding request: %v", err)
    }

    helloResp, err := sayHelloHandler(helloReq)
    if err != nil {
        log.Fatalf("Handler error: %v", err)
    }

    response, err := encodeResponse(helloResp)
    if err != nil {
        log.Fatalf("Error encoding response: %v", err)
    }

    fmt.Println(response)
}
