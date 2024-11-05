package client

import (
    "flag"
    "testing"

    pb "main/genproto"
)

var (
    contextKey = flag.String("contextKey", "camera", "Context key for ad request")
)

func TestGetAds(t *testing.T) {
    flag.Parse()
    adResponse, err := GetAds(
        &pb.AdRequest{ContextKeys: []string{*contextKey}},
        nil,
    )
    if err != nil {
        t.Fatal(err)
    }
    for _, ad := range adResponse.GetAds() {
        t.Logf("Ad: %s", ad.GetText())
    }
    t.Log("Ad retrieval complete.")
}
