// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "github.com/sirupsen/logrus"
    "google.golang.org/protobuf/proto"
    pb "main/genproto"
)

var (
    catalogMutex *sync.Mutex
    log          *logrus.Logger
    extraLatency time.Duration

    svc = &productCatalog{}

    reloadCatalog bool
)

func init() {
    log = logrus.New()
    log.Formatter = &logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "severity",
            logrus.FieldKeyMsg:   "message",
        },
        TimestampFormat: time.RFC3339Nano,
    }
    log.Out = os.Stdout
    catalogMutex = &sync.Mutex{}
}

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

// handleRequest chooses the correct handler function to call
func handleRequest(msg proto.Message, reqData *RequestData) (proto.Message, error) {
    switch reqData.RequestContext.HTTP.Path {
    case "/product-catalog-service/list-products":
        return svc.ListProducts(msg.(*pb.Empty))
    case "/product-catalog-service/get-product":
        return svc.GetProduct(msg.(*pb.GetProductRequest))
    case "/product-catalog-service/search-products":
        return svc.SearchProducts(msg.(*pb.SearchProductsRequest))
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
    case "/product-catalog-service/list-products":
        msg = &pb.Empty{}
    case "/product-catalog-service/get-product":
        msg = &pb.GetProductRequest{}
    case "/product-catalog-service/search-products":
        msg = &pb.SearchProductsRequest{}
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
    flag.Parse()

    // set injected latency
    if s := os.Getenv("EXTRA_LATENCY"); s != "" {
        v, err := time.ParseDuration(s)
        if err != nil {
            log.Fatalf("failed to parse EXTRA_LATENCY (%s) as time.Duration: %+v", v, err)
        }
        extraLatency = v
        log.Infof("extra latency enabled (duration: %v)", extraLatency)
    } else {
        extraLatency = time.Duration(0)
    }

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
    go func() {
        for {
            sig := <-sigs
            log.Printf("Received signal: %s", sig)
            if sig == syscall.SIGUSR1 {
                reloadCatalog = true
                log.Infof("Enable catalog reloading")
            } else {
                reloadCatalog = false
                log.Infof("Disable catalog reloading")
            }
        }
    }()

    reader := bufio.NewReader(os.Stdin)
    request, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalf("Failed to read from stdin: %v", err)
    }
    request = request[:len(request)-1]

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

    select {}
}
