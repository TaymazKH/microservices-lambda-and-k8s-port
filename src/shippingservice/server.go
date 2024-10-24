package main

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/proto"

    pb "main/genproto"
)

var runningInLambda = os.Getenv("RUN_LAMBDA") == "1"

const (
    defaultPort = "8080"

    getQuoteRPC  = "get-quote"
    shipOrderRPC = "ship-order"
)

// callRPC chooses the correct handler function to call.
func callRPC(msg *proto.Message, reqData *RequestData) (proto.Message, error) {
    switch reqData.Headers["rpc-name"] {
    case getQuoteRPC:
        return handleGetQuote((*msg).(*pb.GetQuoteRequest), &reqData.Headers)
    default:
        return handleShipOrder((*msg).(*pb.ShipOrderRequest), &reqData.Headers)
    }
}

// determineMessageType chooses the correct message type to initialize.
func determineMessageType(rpcName string) proto.Message {
    switch rpcName {
    case getQuoteRPC:
        return &pb.GetQuoteRequest{}
    case shipOrderRPC:
        return &pb.ShipOrderRequest{}
    default:
        return nil
    }
}

// RequestData represents the structure of the incoming JSON string or HTTP request.
type RequestData struct {
    Body            string            `json:"body"`
    Headers         map[string]string `json:"headers"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
    BinBody         []byte
}

// ResponseData represents the structure of the outgoing JSON string or HTTP request.
type ResponseData struct {
    StatusCode      int               `json:"statusCode"`
    Headers         map[string]string `json:"headers"`
    Body            string            `json:"body"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
    BinBody         []byte
}

// decodeRequest decodes the incoming RequestData into a protobuf message.
// returns a ResponseData in case of an invalid request.
func decodeRequest(reqData *RequestData) (*proto.Message, *ResponseData, error) {
    var binReqBody []byte
    if reqData.IsBase64Encoded {
        var err error
        binReqBody, err = base64.StdEncoding.DecodeString(reqData.Body)
        if err != nil {
            return nil, nil, fmt.Errorf("failed to decode base64 body: %w", err)
        }
    } else {
        binReqBody = reqData.BinBody
    }

    rpcName := reqData.Headers["rpc-name"]
    msg := determineMessageType(rpcName)
    if msg == nil {
        return nil, generateErrorResponse(codes.Unimplemented, fmt.Sprintf("unknown RPC name: %s", rpcName)), nil
    }

    if err := proto.Unmarshal(binReqBody, msg); err != nil {
        return nil, generateErrorResponse(codes.InvalidArgument, err.Error()), nil
    }

    return &msg, nil, nil
}

// encodeResponse encodes a protobuf response message or an error into a ResponseData.
func encodeResponse(msg *proto.Message, rpcError error) (*ResponseData, error) {
    if rpcError != nil {
        stat := status.Convert(rpcError)
        return generateErrorResponse(stat.Code(), stat.Message()), nil
    }

    binRespBody, err := proto.Marshal(*msg)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal response: %w", err)
    }

    if runningInLambda {
        return &ResponseData{
            StatusCode: 200,
            Headers: map[string]string{
                "content-type": "application/octet-stream",
                "grpc-status":  strconv.Itoa(int(codes.OK))},
            Body:            base64.StdEncoding.EncodeToString(binRespBody),
            IsBase64Encoded: true,
        }, nil

    } else {
        return &ResponseData{
            StatusCode: 200,
            Headers: map[string]string{
                "content-type": "application/octet-stream",
                "grpc-status":  strconv.Itoa(int(codes.OK))},
            BinBody:         binRespBody,
            IsBase64Encoded: false,
        }, nil
    }
}

// generateErrorResponse creates a ResponseData from a message and gRPC status code.
func generateErrorResponse(code codes.Code, message string) *ResponseData {
    return &ResponseData{
        StatusCode: 200,
        Headers: map[string]string{
            "content-type": "text/plain",
            "grpc-status":  strconv.Itoa(int(code))},
        Body:            message,
        IsBase64Encoded: false,
    }
}

func runLambda() error {
    reader := bufio.NewReader(os.Stdin)
    request, err := reader.ReadString('\n')
    if err != nil {
        return fmt.Errorf("failed to read from stdin: %w", err)
    }
    request = strings.TrimSpace(request)

    var reqData *RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        return fmt.Errorf("failed to parse request JSON: %w", err)
    }

    var respData *ResponseData
    reqMsg, respData, err := decodeRequest(reqData)
    if err != nil {
        return fmt.Errorf("error decoding request: %w", err)

    } else if respData == nil {
        respMsg, rpcError := callRPC(reqMsg, reqData)

        respData, err = encodeResponse(&respMsg, rpcError)
        if err != nil {
            return fmt.Errorf("error encoding response: %w", err)
        }
    }

    jsonResponse, err := json.Marshal(respData)
    if err != nil {
        return fmt.Errorf("failed to marshal JSON response: %w", err)
    }

    fmt.Println(string(jsonResponse))
    return nil
}

func runHTTPServer() error {
    httpHandler := func(w http.ResponseWriter, r *http.Request) {
        reqBody, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Error reading request body: %v", err)
            http.Error(w, "failed to read request body", http.StatusInternalServerError)
            return
        }
        defer r.Body.Close()

        headers := make(map[string]string)
        for k, vs := range r.Header {
            s := vs[0]
            for _, v := range vs[1:] {
                s += "," + v
            }
            headers[strings.ToLower(k)] = s
        }

        reqData := &RequestData{
            BinBody:         reqBody,
            Headers:         headers,
            IsBase64Encoded: false,
        }

        var respData *ResponseData
        reqMsg, respData, err := decodeRequest(reqData)
        if err != nil {
            log.Printf("Error decoding request: %v", err)
            http.Error(w, "failed to decode request", http.StatusInternalServerError)
            return

        } else if respData == nil {
            respMsg, rpcError := callRPC(reqMsg, reqData)

            respData, err = encodeResponse(&respMsg, rpcError)
            if err != nil {
                log.Printf("Error encoding response: %v", err)
                http.Error(w, "failed to encode response", http.StatusInternalServerError)
                return
            }
        }

        for k, v := range respData.Headers {
            w.Header().Set(k, v)
        }
        w.WriteHeader(respData.StatusCode)
        if _, err := w.Write(respData.BinBody); err != nil {
            log.Printf("Error writing response: %v", err)
        }
    }

    port := defaultPort
    if p, ok := os.LookupEnv("PORT"); ok {
        port = p
    }
    log.Println("Port:", port)

    http.HandleFunc("/", httpHandler)
    return http.ListenAndServe(":"+port, nil)
}

func main() {
    if runningInLambda {
        log.Println("Running Lambda handler.")
        if err := runLambda(); err != nil {
            log.Fatalf("Error running lambda handler: %v", err)
        }
    } else {
        log.Println("Running HTTP server.")
        if err := runHTTPServer(); err != nil {
            log.Fatalf("HTTP server ended with error: %v", err)
        }
    }
}
