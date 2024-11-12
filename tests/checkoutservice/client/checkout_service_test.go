package client

import (
    "fmt"
    "testing"

    "google.golang.org/protobuf/encoding/protojson"

    pb "main/genproto"
)

func TestPlaceOrder(t *testing.T) {
    if err := createFakeCart(); err != nil {
        t.Fatalf("Populating cart failed: %v", err)
    }

    request := &pb.PlaceOrderRequest{
        UserId:       "user0",
        UserCurrency: "USD",
        Address: &pb.Address{
            Country:       "USA",
            State:         "CA",
            City:          "city0",
            StreetAddress: "street0",
            ZipCode:       95000,
        },
        Email: "user0@example.com",
        CreditCard: &pb.CreditCardInfo{
            CreditCardNumber:          "4111111111111111",
            CreditCardCvv:             123,
            CreditCardExpirationYear:  3000,
            CreditCardExpirationMonth: 1,
        },
    }

    response, err := PlaceOrder(request, nil)
    if err != nil {
        t.Fatalf("PlaceOrder failed: %v", err)
    }

    jsonResponse, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(response)
    if err != nil {
        t.Fatalf("Marshaling response failed: %v", err)
    }

    t.Logf("PlaceOrder response:\n%v", string(jsonResponse))
}

func createFakeCart() error {
    _, err := EmptyCart(
        &pb.EmptyCartRequest{UserId: "user0"},
        nil,
    )
    if err != nil {
        return fmt.Errorf("empty cart failed: %v", err)
    }

    _, err = AddItem(
        &pb.AddItemRequest{UserId: "user0", Item: &pb.CartItem{ProductId: "OLJCESPC7Z", Quantity: 2}}, // sunglasses
        nil,
    )
    if err != nil {
        return fmt.Errorf("add item failed: %v", err)
    }

    return nil
}
