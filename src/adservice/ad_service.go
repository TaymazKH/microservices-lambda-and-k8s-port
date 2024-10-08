package main

import (
    "log"
    "math/rand"
    "time"

    pb "main/genproto"
)

// AdService struct holds the data and methods for the ad service
type AdService struct {
    maxAdsToServe int
    adsMap        map[string][]*pb.Ad
}

func NewAdService() *AdService {
    return &AdService{
        maxAdsToServe: 2,
        adsMap:        createAdsMap(),
    }
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

// GetAds processes the AdRequest and returns an AdResponse
func (s *AdService) GetAds(req *pb.AdRequest, headers *map[string]string) (*pb.AdResponse, error) {
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

    return &pb.AdResponse{Ads: allAds}, nil
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
