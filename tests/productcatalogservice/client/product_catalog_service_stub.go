package client

import (
    "log"
    "net/http"
    "os"
    "strconv"

    pb "main/genproto"
)

var (
    productCatalogServiceAddr    *string
    productCatalogServiceTimeout *int
)

const (
    productCatalogService = "product-catalog-service"
    listProductsRPC       = "list-products"
    getProductRPC         = "get-product"
    searchProductsRPC     = "search-products"
)

// ListProducts represents the ProductCatalogService/ListProducts RPC.
// context can be sent as custom headers.
func ListProducts(request *pb.Empty, header *http.Header) (*pb.ListProductsResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*productCatalogServiceAddr, productCatalogService, listProductsRPC, &binReq, header, *productCatalogServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, listProductsRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.ListProductsResponse), nil
}

// GetProduct represents the ProductCatalogService/GetProduct RPC.
// context can be sent as custom headers.
func GetProduct(request *pb.GetProductRequest, header *http.Header) (*pb.Product, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*productCatalogServiceAddr, productCatalogService, getProductRPC, &binReq, header, *productCatalogServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, getProductRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.Product), nil
}

// SearchProducts represents the ProductCatalogService/SearchProducts RPC.
// context can be sent as custom headers.
func SearchProducts(request *pb.SearchProductsRequest, header *http.Header) (*pb.SearchProductsResponse, error) {
    binReq, err := marshalRequest(request)
    if err != nil {
        return nil, err
    }

    respBody, header, err := sendRequest(*productCatalogServiceAddr, productCatalogService, searchProductsRPC, &binReq, header, *productCatalogServiceTimeout)
    if err != nil {
        return nil, err
    }

    msg, err := unmarshalResponse(respBody, header, searchProductsRPC)
    if err != nil {
        return nil, err
    }

    return (*msg).(*pb.SearchProductsResponse), nil
}

// init loads the address and timeout variables.
func init() {
    a, ok := os.LookupEnv("PRODUCT_CATALOG_SERVICE_ADDR")
    if !ok {
        log.Fatal("PRODUCT_CATALOG_SERVICE_ADDR environment variable not set")
    }
    productCatalogServiceAddr = &a

    t, ok := os.LookupEnv("PRODUCT_CATALOG_SERVICE_TIMEOUT")
    if !ok {
        t := defaultTimeout
        productCatalogServiceTimeout = &t
    } else {
        if t, err := strconv.Atoi(t); err != nil || t <= 0 {
            t = defaultTimeout
            productCatalogServiceTimeout = &t
        } else {
            productCatalogServiceTimeout = &t
        }
    }
}
