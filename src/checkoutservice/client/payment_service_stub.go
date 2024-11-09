package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    paymentServiceAddr    *string
    paymentServiceTimeout *int
)

const (
    paymentService = "payment-service"
    chargeRPC      = "charge"
)

// Charge represents the PaymentService/Charge RPC.
// context can be sent as custom headers.
func Charge(request *pb.ChargeRequest, header *http.Header) (*pb.ChargeResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*paymentServiceAddr, paymentService, chargeRPC, &binReq, header, *paymentServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, chargeRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ChargeResponse), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("PAYMENT_SERVICE_ADDR")
    if !ok {
        log.Fatal("PAYMENT_SERVICE_ADDR environment variable not set")
    }
    paymentServiceAddr = &a

    t, ok := os.LookupEnv("PAYMENT_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        paymentServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            paymentServiceTimeout = &t
        } else {
            paymentServiceTimeout = &t
        }
    }
}
