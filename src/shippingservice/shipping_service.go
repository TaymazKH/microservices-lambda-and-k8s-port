package main

import (
    "fmt"
    "log"

    pb "main/genproto"
)

// GetQuote produces a shipping quote (cost) in USD.
func handleGetQuote(in *pb.GetQuoteRequest, headers *map[string]string) (*pb.GetQuoteResponse, error) {
    log.Print("[GetQuote] received request")
    defer log.Print("[GetQuote] completed request")

    // 1. Generate a quote based on the total number of items to be shipped.
    quote := CreateQuoteFromCount(0)

    // 2. Generate a response.
    return &pb.GetQuoteResponse{
        CostUsd: &pb.Money{
            CurrencyCode: "USD",
            Units:        int64(quote.Dollars),
            Nanos:        int32(quote.Cents * 10000000)},
    }, nil
}

// ShipOrder mocks that the requested items will be shipped.
// It supplies a tracking ID for notional lookup of shipment delivery status.
func handleShipOrder(in *pb.ShipOrderRequest, headers *map[string]string) (*pb.ShipOrderResponse, error) {
    log.Print("[ShipOrder] received request")
    defer log.Print("[ShipOrder] completed request")
    // 1. Create a Tracking ID
    baseAddress := fmt.Sprintf("%s, %s, %s", in.Address.StreetAddress, in.Address.City, in.Address.State)
    id := CreateTrackingId(baseAddress)

    // 2. Generate a response.
    return &pb.ShipOrderResponse{
        TrackingId: id,
    }, nil
}
