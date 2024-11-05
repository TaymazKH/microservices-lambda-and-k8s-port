package client

import (
    "flag"
    "testing"

    pb "main/genproto"
)

// TestGetQuote is a basic check on the GetQuote RPC service.
func TestGetQuote(t *testing.T) {
    flag.Parse()

    // A basic test case to test logic and protobuf interactions.
    req := &pb.GetQuoteRequest{
        Address: &pb.Address{
            StreetAddress: "Muffin Man",
            City:          "London",
            State:         "",
            Country:       "England",
        },
        Items: []*pb.CartItem{
            {
                ProductId: "23",
                Quantity:  1,
            },
            {
                ProductId: "46",
                Quantity:  3,
            },
        },
    }

    res, err := GetQuote(req, nil)
    if err != nil {
        t.Errorf("TestGetQuote (%v) failed", err)
    } else if res.CostUsd.GetUnits() != 8 || res.CostUsd.GetNanos() != 990000000 {
        t.Errorf("TestGetQuote: Quote value '%d.%d' does not match expected '%s'", res.CostUsd.GetUnits(), res.CostUsd.GetNanos(), "11.220000000")
    }
}

// TestShipOrder is a basic check on the ShipOrder RPC service.
func TestShipOrder(t *testing.T) {
    flag.Parse()

    // A basic test case to test logic and protobuf interactions.
    req := &pb.ShipOrderRequest{
        Address: &pb.Address{
            StreetAddress: "Muffin Man",
            City:          "London",
            State:         "",
            Country:       "England",
        },
        Items: []*pb.CartItem{
            {
                ProductId: "23",
                Quantity:  1,
            },
            {
                ProductId: "46",
                Quantity:  3,
            },
        },
    }

    res, err := ShipOrder(req, nil)
    if err != nil {
        t.Errorf("TestShipOrder (%v) failed", err)
    } else if len(res.TrackingId) != 18 {
        // @todo improve quality of this test to check for a pattern such as '[A-Z]{2}-\d+-\d+'.
        t.Errorf("TestShipOrder: Tracking ID is malformed - has %d characters, %d expected", len(res.TrackingId), 18)
    }
}
