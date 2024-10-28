package main

import (
    "main/cartstore"
    pb "main/genproto"
)

type CartService struct {
    cartStore cartstore.CartStore
}

func NewCartService(cartStore cartstore.CartStore) *CartService {
    return &CartService{cartStore: cartStore}
}

func (s *CartService) AddItem(req *pb.AddItemRequest, headers *map[string]string) (*pb.Empty, error) {
    err := s.cartStore.AddItemAsync(req.UserId, req.Item.ProductId, req.Item.Quantity)
    if err != nil {
        return nil, err
    }
    return &pb.Empty{}, nil
}

func (s *CartService) GetCart(req *pb.GetCartRequest, headers *map[string]string) (*pb.Cart, error) {
    cart, err := s.cartStore.GetCartAsync(req.UserId)
    if err != nil {
        return nil, err
    }
    return cart, nil
}

func (s *CartService) EmptyCart(req *pb.EmptyCartRequest, headers *map[string]string) (*pb.Empty, error) {
    err := s.cartStore.EmptyCartAsync(req.UserId)
    if err != nil {
        return nil, err
    }
    return &pb.Empty{}, nil
}
