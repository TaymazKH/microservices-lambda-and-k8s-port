package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    adServiceAddr    *string
    adServiceTimeout *int
)

const (
    adService = "ad-service"
    getAdsRPC = "get-ads"
)

// GetAds represents the AdService/GetAds RPC.
// context can be sent as custom headers.
func GetAds(request *pb.AdRequest, header *http.Header) (*pb.AdResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*adServiceAddr, adService, getAdsRPC, &binReq, header, *adServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, getAdsRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.AdResponse), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("AD_SERVICE_ADDR")
    if !ok {
        log.Fatal("AD_SERVICE_ADDR environment variable not set")
    }
    adServiceAddr = &a

    t, ok := os.LookupEnv("AD_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        adServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            adServiceTimeout = &t
        } else {
            adServiceTimeout = &t
        }
    }
}
