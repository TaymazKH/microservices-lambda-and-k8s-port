package client

import (
    "flag"
    "testing"

    pb "main/genproto"
)

func TestCartService(t *testing.T) {
    flag.Parse()

    userId := "user0"
    productId := "product0"

    // Empty the cart
    _, err := EmptyCart(&pb.EmptyCartRequest{UserId: userId}, nil)
    if err != nil {
        t.Fatal(err)
    }

    // Check if cart is empty
    cart, err := GetCart(&pb.GetCartRequest{UserId: userId}, nil)
    if err != nil {
        t.Fatal(err)
    }
    if cart.UserId != userId {
        t.Fatalf("cart user id is %s, expected %s", cart.UserId, userId)
    }
    if len(cart.Items) != 0 {
        t.Fatal("cart is not empty")
    }

    // Add an item
    _, err = AddItem(
        &pb.AddItemRequest{UserId: userId, Item: &pb.CartItem{ProductId: productId, Quantity: 2}},
        nil,
    )
    if err != nil {
        t.Fatal(err)
    }

    // Check if the cart has that one item
    cart, err = GetCart(&pb.GetCartRequest{UserId: userId}, nil)
    if err != nil {
        t.Fatal(err)
    }
    if len(cart.Items) != 1 {
        t.Fatal("cart doesn't have exactly one item")
    }
    if cart.Items[0].ProductId != productId {
        t.Fatalf("product id is %s, expected %s", cart.Items[0].ProductId, productId)
    }
    if cart.Items[0].Quantity != 2 {
        t.Fatalf("quantity is %d, expected %d", cart.Items[0].Quantity, 2)
    }

    // Empty the cart
    _, err = EmptyCart(&pb.EmptyCartRequest{UserId: userId}, nil)
    if err != nil {
        t.Fatal(err)
    }
}
