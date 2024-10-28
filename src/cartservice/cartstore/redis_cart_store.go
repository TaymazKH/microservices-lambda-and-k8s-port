package cartstore

import (
    "context"
    "encoding/json" // todo: use protojson
    "errors"
    "fmt"
    "log"

    "github.com/redis/go-redis/v9"

    pb "main/genproto"
)

type JsonCart struct {
    UserId string         `json:"userId"`
    Items  []JsonCartItem `json:"items"`
}

type JsonCartItem struct {
    ProductId string `json:"productId"`
    Quantity  int    `json:"quantity"`
}

type RedisCartStore struct {
    rdb *redis.Client
    ctx context.Context
}

func NewRedisCartStore(redisAddr, redisPassword string) *RedisCartStore {
    ctx := context.Background()
    rdb := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: redisPassword, // no password set if empty
        DB:       0,             // use default DB
    })
    return &RedisCartStore{rdb: rdb, ctx: ctx}
}

func (store *RedisCartStore) AddItemAsync(userId, productId string, quantity int32) error {
    log.Printf("AddItemAsync called with userId=%s, productId=%s, quantity=%d\n", userId, productId, quantity)

    // Get existing cart from Redis
    val, err := store.rdb.Get(store.ctx, userId).Result()
    cart := &pb.Cart{UserId: userId}

    if errors.Is(err, redis.Nil) {
        // Cart not found, create a new one
        cart.Items = []*pb.CartItem{
            {ProductId: productId, Quantity: quantity},
        }
    } else if err != nil {
        return fmt.Errorf("can't access cart storage: %v", err)
    } else {
        // Unmarshal existing cart and update it
        err = json.Unmarshal([]byte(val), cart) // todo: define and use new struct
        if err != nil {
            return fmt.Errorf("error parsing cart: %v", err)
        }

        // Find existing item in the cart
        itemFound := false
        for i, item := range cart.Items {
            if item.ProductId == productId {
                cart.Items[i].Quantity += quantity
                itemFound = true
                break
            }
        }

        // If item not found, add it
        if !itemFound {
            cart.Items = append(cart.Items, &pb.CartItem{ProductId: productId, Quantity: quantity})
        }
    }

    // Marshal the updated cart and store it in Redis
    cartData, err := json.Marshal(cart) // todo: define and use new struct
    if err != nil {
        return fmt.Errorf("error serializing cart: %v", err)
    }

    err = store.rdb.Set(store.ctx, userId, cartData, 0).Err()
    if err != nil {
        return fmt.Errorf("error storing cart in Redis: %v", err)
    }

    return nil
}

func (store *RedisCartStore) GetCartAsync(userId string) (*pb.Cart, error) {
    log.Printf("GetCartAsync called with userId=%s\n", userId)

    // Get the cart data from Redis
    val, err := store.rdb.Get(store.ctx, userId).Result()
    if errors.Is(err, redis.Nil) {
        // Cart not found, return an empty cart
        return &pb.Cart{UserId: userId}, nil
    } else if err != nil {
        return nil, fmt.Errorf("can't access cart storage: %v", err)
    }

    // Parse the cart from Redis
    cart := &pb.Cart{}
    err = json.Unmarshal([]byte(val), cart) // todo: define and use new struct
    if err != nil {
        return nil, fmt.Errorf("error parsing cart: %v", err)
    }

    return cart, nil
}

func (store *RedisCartStore) EmptyCartAsync(userId string) error {
    log.Printf("EmptyCartAsync called with userId=%s\n", userId)

    // Create a new empty cart
    cart := pb.Cart{UserId: userId}

    // Store the empty cart in Redis
    cartData, err := json.Marshal(cart) // todo: define and use new struct
    if err != nil {
        return fmt.Errorf("error serializing empty cart: %v", err)
    }

    err = store.rdb.Set(store.ctx, userId, cartData, 0).Err()
    if err != nil {
        return fmt.Errorf("error storing empty cart in Redis: %v", err)
    }

    return nil
}

func (store *RedisCartStore) Ping() bool {
    pong, err := store.rdb.Ping(store.ctx).Result()
    if err != nil {
        log.Printf("Error pinging Redis: %v\n", err)
        return false
    }
    log.Printf("Redis Ping result: %s\n", pong)
    return true
}
