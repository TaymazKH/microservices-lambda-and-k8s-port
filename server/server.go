package main

import (
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"

    "google.golang.org/protobuf/proto"
    pb "main/hello"
)

var (
    port = flag.Int("port", 50051, "The server port")
)

// sayHelloHandler handles the /sayhello endpoint
func sayHelloHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Read the request body into a byte slice
    reqBody, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    // Unmarshal the request body into a HelloRequest message
    var req pb.HelloRequest
    if err := proto.Unmarshal(reqBody, &req); err != nil {
        http.Error(w, "Failed to parse request", http.StatusBadRequest)
        return
    }

    // Log the received name
    log.Printf("Received: %v", req.GetName())

    // Create a HelloResponse message
    resp := &pb.HelloResponse{
        Text: "Hello " + req.GetName(),
    }

    // Marshal the HelloResponse message into binary format
    respBody, err := proto.Marshal(resp)
    if err != nil {
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }

    // Set the content type and write the response body
    w.Header().Set("Content-Type", "application/octet-stream")
    _, err = w.Write(respBody)
    if err != nil {
        http.Error(w, "Failed to write response", http.StatusInternalServerError)
    }
}

func main() {
    flag.Parse()

    // Set up the HTTP server and route
    http.HandleFunc("/sayhello", sayHelloHandler)

    addr := fmt.Sprintf(":%d", *port)
    log.Printf("Server listening on %v", addr)

    // Start the HTTP server
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
