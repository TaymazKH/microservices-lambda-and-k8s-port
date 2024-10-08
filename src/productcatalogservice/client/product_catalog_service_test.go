package client

import (
    "flag"
    "testing"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    pb "main/genproto"
)

func TestGetProductExists(t *testing.T) {
    flag.Parse()
    product, err := GetProduct(
        &pb.GetProductRequest{Id: "OLJCESPC7Z"},
        nil,
    )
    if err != nil {
        t.Fatal(err)
    }
    if got, want := product.Name, "Sunglasses"; got != want {
        t.Errorf("got %s, want %s", got, want)
    }
}

func TestGetProductNotFound(t *testing.T) {
    flag.Parse()
    _, err := GetProduct(
        &pb.GetProductRequest{Id: "abc005"},
        nil,
    )
    if got, want := status.Code(err), codes.NotFound; got != want {
        t.Errorf("got %s, want %s", got, want)
    }
}

func TestListProducts(t *testing.T) {
    flag.Parse()
    products, err := ListProducts(
        &pb.Empty{},
        nil,
    )
    if err != nil {
        t.Fatal(err)
    }
    if got, want := len(products.Products), 9; got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}

func TestSearchProducts(t *testing.T) {
    flag.Parse()
    products, err := SearchProducts(
        &pb.SearchProductsRequest{Query: "Outfit"},
        nil,
    )
    if err != nil {
        t.Fatal(err)
    }
    if got, want := len(products.Results), 2; got != want {
        t.Errorf("got %d, want %d", got, want)
    }
}
