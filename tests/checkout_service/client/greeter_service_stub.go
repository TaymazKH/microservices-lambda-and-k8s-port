package main

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    checkoutServiceAddr    *string
    checkoutServiceTimeout *int
)

const (
    checkoutService = "checkout-service"
    placeOrderRPC   = "place-order"
)

// PlaceOrder represents the CheckoutService/PlaceOrder RPC.
// context can be sent as custom headers.
func PlaceOrder(request *pb.PlaceOrderRequest, header *http.Header) (*pb.PlaceOrderResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*checkoutServiceAddr, checkoutService, placeOrderRPC, &binReq, header, *checkoutServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, placeOrderRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.PlaceOrderResponse), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("CHECKOUT_SERVICE_ADDR")
    if !ok {
        log.Fatal("CHECKOUT_SERVICE_ADDR environment variable not set")
    }
    checkoutServiceAddr = &a

    t, ok := os.LookupEnv("CHECKOUT_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        checkoutServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            checkoutServiceTimeout = &t
        } else {
            checkoutServiceTimeout = &t
        }
    }
}
