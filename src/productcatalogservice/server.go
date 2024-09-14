package main

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "strconv"
    "sync"
    "time"

    //"os/signal"
    //"syscall"

    "github.com/sirupsen/logrus"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
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

const (
    productCatalogService = "product-catalog-service"
    listProductsRPC       = "list-products"
    getProductRPC         = "get-product"
    searchProductsRPC     = "search-products"
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
    log.Out = os.Stderr
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
func handleRequest(msg *proto.Message, reqData *RequestData) (proto.Message, error) {
    switch reqData.RequestContext.HTTP.Path {
    case fmt.Sprintf("/%s/%s", productCatalogService, listProductsRPC):
        return svc.ListProducts((*msg).(*pb.Empty))
    case fmt.Sprintf("/%s/%s", productCatalogService, getProductRPC):
        return svc.GetProduct((*msg).(*pb.GetProductRequest))
    default:
        return svc.SearchProducts((*msg).(*pb.SearchProductsRequest))
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
    case fmt.Sprintf("/%s/%s", productCatalogService, listProductsRPC):
        msg = &pb.Empty{}
    case fmt.Sprintf("/%s/%s", productCatalogService, getProductRPC):
        msg = &pb.GetProductRequest{}
    case fmt.Sprintf("/%s/%s", productCatalogService, searchProductsRPC):
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
func encodeResponse(msg *proto.Message, rpcError error) (string, error) {
    var respData ResponseData

    if rpcError == nil {
        binRespBody, err := proto.Marshal(*msg)
        if err != nil {
            return "", fmt.Errorf("failed to marshal response: %w", err)
        }

        encodedRespBody := base64.StdEncoding.EncodeToString(binRespBody) // Base64 encoding is optional.

        respData = ResponseData{
            StatusCode: 200,
            Headers: map[string]string{
                "Content-Type": "application/octet-stream",
                "Grpc-Code":    strconv.Itoa(int(codes.OK))},
            Body:            encodedRespBody, // Use `binRespBody` if not encoded.
            IsBase64Encoded: true,
        }
    } else {
        stat := status.Convert(rpcError)

        respData = ResponseData{
            StatusCode: 200,
            Headers: map[string]string{
                "Content-Type": "text/plain",
                "Grpc-Code":    strconv.Itoa(int(stat.Code()))},
            Body:            stat.Message(),
            IsBase64Encoded: false,
        }
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

    reloadCatalog = false // disabled reload signaling. todo: perhaps debug and enable later?
    //sigs := make(chan os.Signal, 1)
    //signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
    //go func() {
    //    for {
    //        sig := <-sigs
    //        log.Printf("Received signal: %s", sig)
    //        if sig == syscall.SIGUSR1 {
    //            reloadCatalog = true
    //            log.Infof("Enable catalog reloading")
    //        } else {
    //            reloadCatalog = false
    //            log.Infof("Disable catalog reloading")
    //        }
    //    }
    //}()

    reader := bufio.NewReader(os.Stdin)
    request, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalf("Failed to read from stdin: %v", err)
    }
    request = request[:len(request)-1]

    reqMsg, reqData, err := decodeRequest(request)
    if err != nil {
        log.Fatalf("Error decoding request: %v", err)
    }

    respMsg, err := handleRequest(reqMsg, reqData)

    response, err := encodeResponse(&respMsg, err)
    if err != nil {
        log.Fatalf("Error encoding response: %v", err)
    }

    fmt.Println(response)

    //select {}
}
