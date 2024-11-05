package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    shippingServiceAddr    *string
    shippingServiceTimeout *int
)

const (
    shippingService = "shipping-service"
    getQuoteRPC     = "get-quote"
    shipOrderRPC    = "ship-order"
)

// GetQuote represents the ShippingService/GetQuote RPC.
// context can be sent as custom headers.
func GetQuote(request *pb.GetQuoteRequest, header *http.Header) (*pb.GetQuoteResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*shippingServiceAddr, shippingService, getQuoteRPC, &binReq, header, *shippingServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, getQuoteRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.GetQuoteResponse), nil
}

// ShipOrder represents the ShippingService/ShipOrder RPC.
// context can be sent as custom headers.
func ShipOrder(request *pb.ShipOrderRequest, header *http.Header) (*pb.ShipOrderResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*shippingServiceAddr, shippingService, shipOrderRPC, &binReq, header, *shippingServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, shipOrderRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ShipOrderResponse), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("SHIPPING_SERVICE_ADDR")
    if !ok {
        log.Fatal("SHIPPING_SERVICE_ADDR environment variable not set")
    }
    shippingServiceAddr = &a

    t, ok := os.LookupEnv("SHIPPING_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        shippingServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            shippingServiceTimeout = &t
        } else {
            shippingServiceTimeout = &t
        }
    }
}
