package cartstore

import (
    "log"
    "sync"

    "google.golang.org/protobuf/proto"

    pb "main/genproto"
)

type InMemoryCartStore struct {
    carts sync.Map
}

func NewInMemoryCartStore() *InMemoryCartStore {
    log.Println("Initializing InMemory CartStore")
    return &InMemoryCartStore{}
}

func (store *InMemoryCartStore) AddItemAsync(userId, productId string, quantity int32) error {
    log.Printf("AddItemAsync called with userId=%s, productId=%s, quantity=%d\n", userId, productId, quantity)

    var cart *pb.Cart

    if value, ok := store.carts.Load(userId); ok {
        cart = value.(*pb.Cart)
    } else {
        cart = &pb.Cart{UserId: userId}
    }

    itemFound := false
    for _, item := range cart.Items {
        if item.ProductId == productId {
            item.Quantity += quantity
            itemFound = true
            break
        }
    }

    if !itemFound {
        cart.Items = append(cart.Items, &pb.CartItem{ProductId: productId, Quantity: quantity})
    }

    store.carts.Store(userId, cart)
    return nil
}

func (store *InMemoryCartStore) GetCartAsync(userId string) (*pb.Cart, error) {
    log.Printf("GetCartAsync called with userId=%s\n", userId)

    if value, ok := store.carts.Load(userId); ok {
        return proto.Clone(value.(*pb.Cart)).(*pb.Cart), nil
    }

    return &pb.Cart{UserId: userId}, nil
}

func (store *InMemoryCartStore) EmptyCartAsync(userId string) error {
    log.Printf("EmptyCartAsync called with userId=%s\n", userId)

    store.carts.Delete(userId)
    return nil
}

func (store *InMemoryCartStore) Ping() bool {
    log.Println("InMemory CartStore Ping called - always returns true")
    return true
}
