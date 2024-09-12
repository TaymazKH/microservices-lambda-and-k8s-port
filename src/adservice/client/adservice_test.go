package main

import (
    "flag"
    "fmt"
    "log"
    "testing"

    pb "main/genproto"
)

func TestSearchProducts(t *testing.T) {
    flag.Parse()
    adResponse, err := GetAds(
        &pb.AdRequest{ContextKeys: []string{*contextKey}},
    )
    if err != nil {
        t.Fatal(err)
    }
    for _, ad := range adResponse.GetAds() {
        log.Printf("Ad: %s", ad.GetText())
    }
    fmt.Println("Ad retrieval complete.")
}
