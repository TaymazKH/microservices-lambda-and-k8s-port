package cartstore

import (
    "context"
    "errors"
    "log"

    "github.com/redis/go-redis/v9"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/encoding/protojson"

    pb "main/genproto"
)

type RedisCartStore struct {
    rdb *redis.Client
    ctx context.Context
}

func NewRedisCartStore(redisAddr, redisPassword string) *RedisCartStore {
    log.Printf("Initializing Redis CartStore with address %s", redisAddr)
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

    val, err := store.rdb.Get(store.ctx, userId).Result()
    cart := &pb.Cart{UserId: userId}

    if errors.Is(err, redis.Nil) {
        cart.Items = []*pb.CartItem{
            {ProductId: productId, Quantity: quantity},
        }
    } else if err != nil {
        return status.Errorf(codes.Unavailable, "can't access cart storage: %v", err)
    } else {
        err = protojson.Unmarshal([]byte(val), cart)
        if err != nil {
            return status.Errorf(codes.Internal, "error parsing cart: %v", err)
        }

        itemFound := false
        for i, item := range cart.Items {
            if item.ProductId == productId {
                cart.Items[i].Quantity += quantity
                itemFound = true
                break
            }
        }

        if !itemFound {
            cart.Items = append(cart.Items, &pb.CartItem{ProductId: productId, Quantity: quantity})
        }
    }

    cartData, err := protojson.Marshal(cart)
    if err != nil {
        return status.Errorf(codes.Internal, "error serializing cart: %v", err)
    }

    err = store.rdb.Set(store.ctx, userId, cartData, 0).Err()
    if err != nil {
        return status.Errorf(codes.Unavailable, "error storing cart in Redis: %v", err)
    }

    return nil
}

func (store *RedisCartStore) GetCartAsync(userId string) (*pb.Cart, error) {
    log.Printf("GetCartAsync called with userId=%s\n", userId)

    val, err := store.rdb.Get(store.ctx, userId).Result()
    if errors.Is(err, redis.Nil) {
        return &pb.Cart{UserId: userId}, nil
    } else if err != nil {
        return nil, status.Errorf(codes.Unavailable, "can't access cart storage: %v", err)
    }

    cart := &pb.Cart{}
    err = protojson.Unmarshal([]byte(val), cart)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "error parsing cart: %v", err)
    }

    return cart, nil
}

func (store *RedisCartStore) EmptyCartAsync(userId string) error {
    log.Printf("EmptyCartAsync called with userId=%s\n", userId)

    cart := &pb.Cart{UserId: userId}

    cartData, err := protojson.Marshal(cart)
    if err != nil {
        return status.Errorf(codes.Internal, "error serializing empty cart: %v", err)
    }

    err = store.rdb.Set(store.ctx, userId, cartData, 0).Err()
    if err != nil {
        return status.Errorf(codes.Unavailable, "error storing empty cart in Redis: %v", err)
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
