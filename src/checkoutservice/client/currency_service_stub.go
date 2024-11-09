package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    currencyServiceAddr    *string
    currencyServiceTimeout *int
)

const (
    currencyService           = "currency-service"
    getSupportedCurrenciesRPC = "get-supported-currencies"
    convertRPC                = "convert"
)

// GetSupportedCurrencies represents the CurrencyService/GetSupportedCurrencies RPC.
// context can be sent as custom headers.
func GetSupportedCurrencies(empty *pb.Empty, header *http.Header) (*pb.GetSupportedCurrenciesResponse, error) {
    binReq, err := marshalRequest(empty)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*currencyServiceAddr, currencyService, getSupportedCurrenciesRPC, &binReq, header, *currencyServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, getSupportedCurrenciesRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.GetSupportedCurrenciesResponse), nil
}

// Convert represents the CurrencyService/Convert RPC.
// context can be sent as custom headers.
func Convert(request *pb.CurrencyConversionRequest, header *http.Header) (*pb.Money, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*currencyServiceAddr, currencyService, convertRPC, &binReq, header, *currencyServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, convertRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.Money), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("CURRENCY_SERVICE_ADDR")
    if !ok {
        log.Fatal("CURRENCY_SERVICE_ADDR environment variable not set")
    }
    currencyServiceAddr = &a

    t, ok := os.LookupEnv("CURRENCY_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        currencyServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            currencyServiceTimeout = &t
        } else {
            currencyServiceTimeout = &t
        }
    }
}
