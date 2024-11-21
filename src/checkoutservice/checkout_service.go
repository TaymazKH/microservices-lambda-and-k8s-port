package main

import (
    "fmt"
    "os"
    "time"

    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    stubs "main/client"
    pb "main/genproto"
    "main/money"
)

const (
    usdCurrency = "USD"
)

var log *logrus.Logger

func init() {
    log = logrus.New()
    log.Level = logrus.DebugLevel
    log.Formatter = &logrus.JSONFormatter{
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "severity",
            logrus.FieldKeyMsg:   "message",
        },
        TimestampFormat: time.RFC3339Nano,
    }
    log.Out = os.Stdout
}

type checkoutService struct{}

func (cs *checkoutService) PlaceOrder(req *pb.PlaceOrderRequest, headers *map[string]string) (*pb.PlaceOrderResponse, error) {
    log.Infof("[PlaceOrder] user_id=%q user_currency=%q", req.UserId, req.UserCurrency)

    orderID, err := uuid.NewUUID()
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to generate order uuid")
    }

    prep, err := cs.prepareOrderItemsAndShippingQuoteFromCart(req.UserId, req.UserCurrency, req.Address)
    if err != nil {
        return nil, status.Errorf(codes.Internal, err.Error())
    }

    total := pb.Money{CurrencyCode: req.UserCurrency,
        Units: 0,
        Nanos: 0}
    total = money.Must(money.Sum(total, *prep.shippingCostLocalized))
    for _, it := range prep.orderItems {
        multPrice := money.MultiplySlow(*it.Cost, uint32(it.GetItem().GetQuantity()))
        total = money.Must(money.Sum(total, multPrice))
    }

    txID, err := cs.chargeCard(&total, req.CreditCard)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
    }
    log.Infof("payment went through (transaction_id: %s)", txID)

    shippingTrackingID, err := cs.shipOrder(req.Address, prep.cartItems)
    if err != nil {
        return nil, status.Errorf(codes.Unavailable, "shipping error: %+v", err)
    }

    _ = cs.emptyUserCart(req.UserId)

    orderResult := &pb.OrderResult{
        OrderId:            orderID.String(),
        ShippingTrackingId: shippingTrackingID,
        ShippingCost:       prep.shippingCostLocalized,
        ShippingAddress:    req.Address,
        Items:              prep.orderItems,
    }

    if err := cs.sendOrderConfirmation(req.Email, orderResult); err != nil {
        log.Warnf("failed to send order confirmation to %q: %+v", req.Email, err)
    } else {
        log.Infof("order confirmation email sent to %q", req.Email)
    }
    resp := &pb.PlaceOrderResponse{Order: orderResult}
    return resp, nil
}

type orderPrep struct {
    orderItems            []*pb.OrderItem
    cartItems             []*pb.CartItem
    shippingCostLocalized *pb.Money
}

func (cs *checkoutService) prepareOrderItemsAndShippingQuoteFromCart(userID, userCurrency string, address *pb.Address) (orderPrep, error) {
    var out orderPrep
    cartItems, err := cs.getUserCart(userID)
    if err != nil {
        return out, fmt.Errorf("cart failure: %+v", err)
    }
    orderItems, err := cs.prepOrderItems(cartItems, userCurrency)
    if err != nil {
        return out, fmt.Errorf("failed to prepare order: %+v", err)
    }
    shippingUSD, err := cs.quoteShipping(address, cartItems)
    if err != nil {
        return out, fmt.Errorf("shipping quote failure: %+v", err)
    }
    shippingPrice, err := cs.convertCurrency(shippingUSD, userCurrency)
    if err != nil {
        return out, fmt.Errorf("failed to convert shipping cost to currency: %+v", err)
    }

    out.shippingCostLocalized = shippingPrice
    out.cartItems = cartItems
    out.orderItems = orderItems
    return out, nil
}

func (cs *checkoutService) quoteShipping(address *pb.Address, items []*pb.CartItem) (*pb.Money, error) {
    shippingQuote, err := stubs.GetQuote(&pb.GetQuoteRequest{Address: address, Items: items}, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get shipping quote: %+v", err)
    }
    return shippingQuote.GetCostUsd(), nil
}

func (cs *checkoutService) getUserCart(userID string) ([]*pb.CartItem, error) {
    cart, err := stubs.GetCart(&pb.GetCartRequest{UserId: userID}, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get user cart during checkout: %+v", err)
    }
    return cart.GetItems(), nil
}

func (cs *checkoutService) emptyUserCart(userID string) error {
    if _, err := stubs.EmptyCart(&pb.EmptyCartRequest{UserId: userID}, nil); err != nil {
        return fmt.Errorf("failed to empty user cart during checkout: %+v", err)
    }
    return nil
}

func (cs *checkoutService) prepOrderItems(items []*pb.CartItem, userCurrency string) ([]*pb.OrderItem, error) {
    out := make([]*pb.OrderItem, len(items))

    for i, item := range items {
        product, err := stubs.GetProduct(&pb.GetProductRequest{Id: item.GetProductId()}, nil)
        if err != nil {
            return nil, fmt.Errorf("failed to get product #%q", item.GetProductId())
        }
        price, err := cs.convertCurrency(product.GetPriceUsd(), userCurrency)
        if err != nil {
            return nil, fmt.Errorf("failed to convert price of %q to %s", item.GetProductId(), userCurrency)
        }
        out[i] = &pb.OrderItem{
            Item: item,
            Cost: price}
    }
    return out, nil
}

func (cs *checkoutService) convertCurrency(from *pb.Money, toCurrency string) (*pb.Money, error) {
    result, err := stubs.Convert(&pb.CurrencyConversionRequest{From: from, ToCode: toCurrency}, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to convert currency: %+v", err)
    }
    return result, err
}

func (cs *checkoutService) chargeCard(amount *pb.Money, paymentInfo *pb.CreditCardInfo) (string, error) {
    paymentResp, err := stubs.Charge(&pb.ChargeRequest{Amount: amount, CreditCard: paymentInfo}, nil)
    if err != nil {
        return "", fmt.Errorf("could not charge the card: %+v", err)
    }
    return paymentResp.GetTransactionId(), nil
}

func (cs *checkoutService) sendOrderConfirmation(email string, order *pb.OrderResult) error {
    _, err := stubs.SendOrderConfirmation(&pb.SendOrderConfirmationRequest{Email: email, Order: order}, nil)
    return err
}

func (cs *checkoutService) shipOrder(address *pb.Address, items []*pb.CartItem) (string, error) {
    resp, err := stubs.ShipOrder(&pb.ShipOrderRequest{Address: address, Items: items}, nil)
    if err != nil {
        return "", fmt.Errorf("shipment failed: %+v", err)
    }
    return resp.GetTrackingId(), nil
}
