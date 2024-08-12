package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "log"
    "os"

    "google.golang.org/protobuf/proto"
    pb "main/hello"
)

// RequestData represents the structure of the incoming JSON string
type RequestData struct {
    Body     string            `json:"body"`
    Headers  map[string]string `json:"headers"`
    RouteKey string            `json:"routeKey"`
}

func sayHelloHandler(request string) string {
    // Parse the incoming JSON string to RequestData struct
    var reqData RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        log.Fatalf("Failed to parse request JSON: %v", err)
    }

    // Convert the request body string to a byte slice
    binReqBody := []byte(reqData.Body)

    log.Println("Request Body:", string(binReqBody)) // todo

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

    // Convert the binary response to a string and return it
    return string(binRespBody)
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
