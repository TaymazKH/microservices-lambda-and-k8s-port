package client

import (
    "flag"
    "testing"

    pb "main/genproto"
)

func TestGetAds(t *testing.T) {
    flag.Parse()
    adResponse, err := GetAds(
        &pb.AdRequest{ContextKeys: []string{*contextKey}},
    )
    if err != nil {
        t.Fatal(err)
    }
    for _, ad := range adResponse.GetAds() {
        t.Logf("Ad: %s", ad.GetText())
    }
    t.Log("Ad retrieval complete.")
}
