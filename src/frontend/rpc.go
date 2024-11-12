package main

import (
    "context"
    "time"

    "github.com/pkg/errors"

    stubs "main/client"
    pb "main/genproto"
)

const (
    avoidNoopCurrencyConversionRPC = false
)

func (fe *frontendServer) getCurrencies(ctx context.Context) ([]string, error) {
    currs, err := stubs.GetSupportedCurrencies(&pb.Empty{}, nil)
    if err != nil {
        return nil, err
    }
    var out []string
    for _, c := range currs.CurrencyCodes {
        if _, ok := whitelistedCurrencies[c]; ok {
            out = append(out, c)
        }
    }
    return out, nil
}

func (fe *frontendServer) getProducts(ctx context.Context) ([]*pb.Product, error) {
    resp, err := stubs.ListProducts(&pb.Empty{}, nil)
    return resp.GetProducts(), err
}

func (fe *frontendServer) getProduct(ctx context.Context, id string) (*pb.Product, error) {
    resp, err := stubs.GetProduct(&pb.GetProductRequest{Id: id}, nil)
    return resp, err
}

func (fe *frontendServer) getCart(ctx context.Context, userID string) ([]*pb.CartItem, error) {
    resp, err := stubs.GetCart(&pb.GetCartRequest{UserId: userID}, nil)
    return resp.GetItems(), err
}

func (fe *frontendServer) emptyCart(ctx context.Context, userID string) error {
    _, err := stubs.EmptyCart(&pb.EmptyCartRequest{UserId: userID}, nil)
    return err
}

func (fe *frontendServer) insertCart(ctx context.Context, userID, productID string, quantity int32) error {
    _, err := stubs.AddItem(&pb.AddItemRequest{
        UserId: userID,
        Item: &pb.CartItem{
            ProductId: productID,
            Quantity:  quantity},
    }, nil)
    return err
}

func (fe *frontendServer) convertCurrency(ctx context.Context, money *pb.Money, currency string) (*pb.Money, error) {
    if avoidNoopCurrencyConversionRPC && money.GetCurrencyCode() == currency {
        return money, nil
    }
    return stubs.Convert(&pb.CurrencyConversionRequest{From: money, ToCode: currency}, nil)
}

func (fe *frontendServer) getShippingQuote(ctx context.Context, items []*pb.CartItem, currency string) (*pb.Money, error) {
    quote, err := stubs.GetQuote(
        &pb.GetQuoteRequest{
            Address: nil,
            Items:   items},
        nil)
    if err != nil {
        return nil, err
    }
    localized, err := fe.convertCurrency(ctx, quote.GetCostUsd(), currency)
    return localized, errors.Wrap(err, "failed to convert currency for shipping cost")
}

func (fe *frontendServer) getRecommendations(ctx context.Context, userID string, productIDs []string) ([]*pb.Product, error) {
    resp, err := stubs.ListRecommendations(&pb.ListRecommendationsRequest{UserId: userID, ProductIds: productIDs}, nil)
    if err != nil {
        return nil, err
    }
    out := make([]*pb.Product, len(resp.GetProductIds()))
    for i, v := range resp.GetProductIds() {
        p, err := fe.getProduct(ctx, v)
        if err != nil {
            return nil, errors.Wrapf(err, "failed to get recommended product info (#%s)", v)
        }
        out[i] = p
    }
    if len(out) > 4 {
        out = out[:4] // take only first four to fit the UI
    }
    return out, err
}

func (fe *frontendServer) getAd(ctx context.Context, ctxKeys []string) ([]*pb.Ad, error) {
    ctx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
    defer cancel()

    resp, err := stubs.GetAds(&pb.AdRequest{
        ContextKeys: ctxKeys,
    }, nil)
    return resp.GetAds(), errors.Wrap(err, "failed to get ads")
}
