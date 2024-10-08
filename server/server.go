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

const (
    defaultPort = "8080"

    sayHelloRPC = "say-hello"
    sayByeRPC   = "say-bye"
)

// callRPC chooses the correct handler function to call.
func callRPC(msg *proto.Message, reqData *RequestData) (proto.Message, error) {
    switch reqData.Headers["rpc-name"] {
    case sayHelloRPC:
        return handleSayHello((*msg).(*pb.HelloRequest), &reqData.Headers)
    default:
        return handleSayBye((*msg).(*pb.ByeRequest), &reqData.Headers)
    }
}

// determineMessageType chooses the correct message type to initialize.
func determineMessageType(rpcName string) (proto.Message, error) {
    var msg proto.Message
    switch rpcName {
    case sayHelloRPC:
        msg = &pb.HelloRequest{}
    case sayByeRPC:
        msg = &pb.ByeRequest{}
    default:
        return nil, fmt.Errorf("unknown RPC name: %s", rpcName)
    }
    return msg, nil
}

// RequestContext represents the nested context of the request.
type RequestContext struct {
    HTTP struct {
        Method string `json:"method"`
        Path   string `json:"path"`
    } `json:"http"`
}

// RequestData represents the structure of the incoming JSON string or HTTP request.
type RequestData struct {
    Body            string            `json:"body"`
    Headers         map[string]string `json:"headers"`
    RequestContext  RequestContext    `json:"requestContext"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// ResponseData represents the structure of the outgoing JSON string or HTTP request.
type ResponseData struct {
    StatusCode      int               `json:"statusCode"`
    Headers         map[string]string `json:"headers"`
    Body            string            `json:"body"`
    IsBase64Encoded bool              `json:"isBase64Encoded"`
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
        binReqBody = []byte(reqData.Body)
    }

    msg, err := determineMessageType(reqData.Headers["rpc-name"])
    if err != nil {
        return nil, generateErrorResponse(codes.Unimplemented, err.Error()), nil
    }

    if err := proto.Unmarshal(binReqBody, msg); err != nil {
        return nil, generateErrorResponse(codes.InvalidArgument, err.Error()), nil
    }

    return &msg, nil, nil
}

// encodeResponse encodes a protobuf response message or an error into a ResponseData.
func encodeResponse(msg *proto.Message, rpcError error) (*ResponseData, error) {
    var respData *ResponseData

    if rpcError == nil {
        binRespBody, err := proto.Marshal(*msg)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal response: %w", err)
        }

        //encodedRespBody := base64.StdEncoding.EncodeToString(binRespBody) // Base64 encoding is optional in Lambda.

        respData = &ResponseData{
            StatusCode: 200,
            Headers: map[string]string{
                "content-type": "application/octet-stream",
                "grpc-status":  strconv.Itoa(int(codes.OK))},
            //Body:            encodedRespBody, // Use if encoded in base64.
            //IsBase64Encoded: true,
            Body:            string(binRespBody), // Use if not encoded.
            IsBase64Encoded: false,
        }
    } else {
        stat := status.Convert(rpcError)

        respData = generateErrorResponse(stat.Code(), stat.Message())
    }

    return respData, nil
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
    request = request[:len(request)-1] // Trim any trailing newline characters

    var reqData *RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        return fmt.Errorf("failed to parse request JSON: %w", err)
    }

    var respData *ResponseData
    reqMsg, respData, err := decodeRequest(reqData)
    if err != nil {
        return fmt.Errorf("error decoding request: %w", err)

    } else if respData == nil {
        respMsg, err := callRPC(reqMsg, reqData)

        respData, err = encodeResponse(&respMsg, err)
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
        fmt.Println(headers)

        reqContext := RequestContext{}
        reqContext.HTTP.Method = r.Method
        reqContext.HTTP.Path = r.URL.Path

        reqData := &RequestData{
            Body:            string(reqBody),
            Headers:         headers,
            RequestContext:  reqContext,
            IsBase64Encoded: false,
        }

        var respData *ResponseData
        reqMsg, respData, err := decodeRequest(reqData)
        if err != nil {
            log.Printf("Error decoding request: %v", err)
            http.Error(w, "failed to decode request", http.StatusInternalServerError)
            return

        } else if respData == nil {
            respMsg, err := callRPC(reqMsg, reqData)

            respData, err = encodeResponse(&respMsg, err)
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
        if _, err := w.Write([]byte(respData.Body)); err != nil {
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
    if os.Getenv("RUN_LAMBDA") == "1" {
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
