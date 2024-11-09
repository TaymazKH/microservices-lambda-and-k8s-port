package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    emailServiceAddr    *string
    emailServiceTimeout *int
)

const (
    emailService             = "email-service"
    sendOrderConfirmationRPC = "send-order-confirmation"
)

// SendOrderConfirmation represents the EmailService/SendOrderConfirmation RPC.
// context can be sent as custom headers.
func SendOrderConfirmation(request *pb.SendOrderConfirmationRequest, header *http.Header) (*pb.Empty, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*emailServiceAddr, emailService, sendOrderConfirmationRPC, &binReq, header, *emailServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, sendOrderConfirmationRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.Empty), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("EMAIL_SERVICE_ADDR")
    if !ok {
        log.Fatal("EMAIL_SERVICE_ADDR environment variable not set")
    }
    emailServiceAddr = &a

    t, ok := os.LookupEnv("EMAIL_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        emailServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            emailServiceTimeout = &t
        } else {
            emailServiceTimeout = &t
        }
    }
}
