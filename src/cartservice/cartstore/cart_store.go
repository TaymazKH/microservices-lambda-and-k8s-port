package cartstore

import (
    pb "main/genproto"
)

type CartStore interface {
    AddItemAsync(userId, productId string, quantity int) error
    GetCartAsync(userId string) (*pb.Cart, error)
    EmptyCartAsync(userId string) error
}
