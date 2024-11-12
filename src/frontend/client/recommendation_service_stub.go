package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    recommendationServiceAddr    *string
    recommendationServiceTimeout *int
)

const (
    recommendationService  = "recommendation-service"
    listRecommendationsRPC = "list-recommendations"
)

// ListRecommendations represents the RecommendationService/ListRecommendations RPC.
// context can be sent as custom headers.
func ListRecommendations(request *pb.ListRecommendationsRequest, header *http.Header) (*pb.ListRecommendationsResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*recommendationServiceAddr, recommendationService, listRecommendationsRPC, &binReq, header, *recommendationServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, listRecommendationsRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ListRecommendationsResponse), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("RECOMMENDATION_SERVICE_ADDR")
    if !ok {
        log.Fatal("RECOMMENDATION_SERVICE_ADDR environment variable not set")
    }
    recommendationServiceAddr = &a

    t, ok := os.LookupEnv("RECOMMENDATION_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        recommendationServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            recommendationServiceTimeout = &t
        } else {
            recommendationServiceTimeout = &t
        }
    }
}
