package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    cartServiceAddr    *string
    cartServiceTimeout *int
)

const (
    cartService  = "cart-service"
    addItemRPC   = "add-item"
    getCartRPC   = "get-cart"
    emptyCartRPC = "empty-cart"
)

// AddItem represents the CartService/AddItem RPC.
// context can be sent as custom headers.
func AddItem(request *pb.AddItemRequest, header *http.Header) (*pb.Empty, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*cartServiceAddr, cartService, addItemRPC, &binReq, header, *cartServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, addItemRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.Empty), nil
}

// GetCart represents the CartService/GetCart RPC.
// context can be sent as custom headers.
func GetCart(request *pb.GetCartRequest, header *http.Header) (*pb.Cart, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*cartServiceAddr, cartService, getCartRPC, &binReq, header, *cartServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, getCartRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.Cart), nil
}

// EmptyCart represents the CartService/EmptyCart RPC.
// context can be sent as custom headers.
func EmptyCart(request *pb.EmptyCartRequest, header *http.Header) (*pb.Empty, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*cartServiceAddr, cartService, emptyCartRPC, &binReq, header, *cartServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, emptyCartRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.Empty), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("CART_SERVICE_ADDR")
    if !ok {
        log.Fatal("CART_SERVICE_ADDR environment variable not set")
    }
    cartServiceAddr = &a

    t, ok := os.LookupEnv("CART_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        cartServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            cartServiceTimeout = &t
        } else {
            cartServiceTimeout = &t
        }
    }
}
