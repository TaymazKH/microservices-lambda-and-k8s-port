package cartstore

import (
    pb "main/genproto"
)

type CartStore interface {
    AddItemAsync(userId, productId string, quantity int32) error
    GetCartAsync(userId string) (*pb.Cart, error)
    EmptyCartAsync(userId string) error
    Ping() bool
}
