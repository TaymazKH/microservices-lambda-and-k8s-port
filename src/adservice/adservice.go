package main

import (
    "bufio"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "os"
    "strings"
    "time"

    "google.golang.org/protobuf/proto"
    pb "main/genproto"
)

// AdService struct holds the data and methods for the ad service
type AdService struct {
    maxAdsToServe int
    adsMap        map[string][]*pb.Ad
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

func NewAdService() *AdService {
    return &AdService{
        maxAdsToServe: 2,
        adsMap:        createAdsMap(),
    }
}

func (s *AdService) getAdsByCategory(category string) []*pb.Ad {
    return s.adsMap[category]
}

func (s *AdService) getRandomAds() []*pb.Ad {
    var ads []*pb.Ad
    var allAds []*pb.Ad
    for _, categoryAds := range s.adsMap {
        allAds = append(allAds, categoryAds...)
    }

    rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < s.maxAdsToServe; i++ {
        ads = append(ads, allAds[rand.Intn(len(allAds))])
    }

    return ads
}

func createAdsMap() map[string][]*pb.Ad {
    return map[string][]*pb.Ad{
        "clothing": {
            &pb.Ad{RedirectUrl: "/product/66VCHSJNUP", Text: "Tank top for sale. 20% off."},
        },
        "accessories": {
            &pb.Ad{RedirectUrl: "/product/1YMWWN1N4O", Text: "Watch for sale. Buy one, get second kit for free"},
        },
        "footwear": {
            &pb.Ad{RedirectUrl: "/product/L9ECAV7KIM", Text: "Loafers for sale. Buy one, get second one for free"},
        },
        "hair": {
            &pb.Ad{RedirectUrl: "/product/2ZYFJ3GM2N", Text: "Hairdryer for sale. 50% off."},
        },
        "decor": {
            &pb.Ad{RedirectUrl: "/product/0PUK6V6EV0", Text: "Candle holder for sale. 30% off."},
        },
        "kitchen": {
            &pb.Ad{RedirectUrl: "/product/9SIQT8TOJO", Text: "Bamboo glass jar for sale. 10% off."},
            &pb.Ad{RedirectUrl: "/product/6E92ZMYYFZ", Text: "Mug for sale. Buy two, get third one for free"},
        },
    }
}

// decodeRequest decodes the incoming JSON request into the protobuf message
func decodeRequest(request string) (*pb.AdRequest, error) {
    var reqData RequestData
    if err := json.Unmarshal([]byte(request), &reqData); err != nil {
        return nil, fmt.Errorf("failed to parse request JSON: %w", err)
    }

    var binReqBody []byte
    if reqData.IsBase64Encoded {
        var err error
        binReqBody, err = base64.StdEncoding.DecodeString(reqData.Body)
        if err != nil {
            return nil, fmt.Errorf("failed to decode base64 body: %w", err)
        }
    } else {
        binReqBody = []byte(reqData.Body)
    }

    var adRequest pb.AdRequest
    if err := proto.Unmarshal(binReqBody, &adRequest); err != nil {
        return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
    }

    return &adRequest, nil
}

// handleGetAds processes the AdRequest and returns an AdResponse
func (s *AdService) handleGetAds(req *pb.AdRequest) *pb.AdResponse {
    var allAds []*pb.Ad
    log.Printf("Received ad request (context_words=%v)", req.ContextKeys)

    if len(req.ContextKeys) > 0 {
        for _, contextKey := range req.ContextKeys {
            ads := s.getAdsByCategory(contextKey)
            allAds = append(allAds, ads...)
        }
    } else {
        allAds = s.getRandomAds()
    }

    if len(allAds) == 0 {
        allAds = s.getRandomAds()
    }

    return &pb.AdResponse{Ads: allAds}
}

// encodeResponse encodes the protobuf response into the outgoing JSON response
func encodeResponse(adResponse *pb.AdResponse) (string, error) {
    binRespBody, err := proto.Marshal(adResponse)
    if err != nil {
        return "", fmt.Errorf("failed to marshal response: %w", err)
    }

    encodedRespBody := base64.StdEncoding.EncodeToString(binRespBody)

    respData := ResponseData{
        StatusCode:      200,
        Headers:         map[string]string{"Content-Type": "application/octet-stream"},
        Body:            encodedRespBody,
        IsBase64Encoded: true,
    }

    jsonResponse, err := json.Marshal(respData)
    if err != nil {
        return "", fmt.Errorf("failed to marshal JSON response: %w", err)
    }

    return string(jsonResponse), nil
}

func main() {
    service := NewAdService()

    reader := bufio.NewReader(os.Stdin)
    request, err := reader.ReadString('\n')
    if err != nil {
        log.Fatalf("Failed to read from stdin: %v", err)
    }
    request = strings.TrimSpace(request)

    adRequest, err := decodeRequest(request)
    if err != nil {
        log.Fatalf("Error decoding request: %v", err)
    }

    adResponse := service.handleGetAds(adRequest)

    response, err := encodeResponse(adResponse)
    if err != nil {
        log.Fatalf("Error encoding response: %v", err)
    }

    fmt.Println(response)
}
