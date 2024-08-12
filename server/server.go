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

// RequestData represents the structure of the incoming JSON string
type RequestData struct {
    Body            string            `json:"body"`
    Headers         map[string]string `json:"headers"`
    RouteKey        string            `json:"routeKey"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// ResponseData represents the structure of the outgoing JSON string
type ResponseData struct {
    Body            string `json:"body"`
    IsBase64Encoded bool   `json:"isBase64Encoded"`
    StatusCode      int    `json:"statusCode"`
}

func sayHelloHandler(request string) string {
    // Parse the incoming JSON string to RequestData struct
    var reqData RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        log.Fatalf("Failed to parse request JSON: %v", err)
    }

    // Decode the request body if it is base64 encoded
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

    // Unmarshal the binary body into a HelloRequest message
    var helloReq pb.HelloRequest
    if err := proto.Unmarshal(binReqBody, &helloReq); err != nil {
        log.Fatalf("Failed to unmarshal request body: %v", err)
    }

    // Log the received name
    log.Printf("Received: %v", helloReq.GetName())

    // Create a HelloResponse message
    helloResp := &pb.HelloResponse{
        Text: "Hello " + helloReq.GetName(),
    }

    // Marshal the HelloResponse message into binary format
    binRespBody, err := proto.Marshal(helloResp)
    if err != nil {
        log.Fatalf("Failed to marshal response: %v", err)
    }

    // Base64 encode the binary response body for safe inclusion in JSON
    encodedRespBody := base64.StdEncoding.EncodeToString(binRespBody) // Base64 encoding is optional.

    // Create the JSON response structure
    respData := ResponseData{
        Body:            encodedRespBody, // Use `binRespBody` if not encoded.
        IsBase64Encoded: true,
        StatusCode:      200,
    }

    // Marshal the response structure into a JSON string
    jsonResponse, err := json.Marshal(respData)
    if err != nil {
        log.Fatalf("Failed to marshal JSON response: %v", err)
    }

    // Convert the binary response to a string and return it
    return string(jsonResponse)
}

func main() {
    // Read the entire input from stdin
    reader := bufio.NewReader(os.Stdin)
    request, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalf("Failed to read from stdin: %v", err)
    }

    // Trim any trailing newline characters from the input
    request = request[:len(request)-1]

    // Call the handler function with the input string
    response := sayHelloHandler(request)

    // Print the response to stdout
    fmt.Println(response)
}
