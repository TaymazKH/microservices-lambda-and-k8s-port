package main

import (
    "strings"
    "time"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    pb "main/genproto"
)

type productCatalog struct {
    catalog pb.ListProductsResponse
}

func (p *productCatalog) ListProducts(empty *pb.Empty, headers *map[string]string) (*pb.ListProductsResponse, error) {
    time.Sleep(extraLatency)

    return &pb.ListProductsResponse{Products: p.parseCatalog()}, nil
}

func (p *productCatalog) GetProduct(req *pb.GetProductRequest, headers *map[string]string) (*pb.Product, error) {
    time.Sleep(extraLatency)

    var found *pb.Product
    for i := 0; i < len(p.parseCatalog()); i++ {
        if req.Id == p.parseCatalog()[i].Id {
            found = p.parseCatalog()[i]
        }
    }

    if found == nil {
        return nil, status.Errorf(codes.NotFound, "no product with ID %s", req.Id)
    }
    return found, nil
}

func (p *productCatalog) SearchProducts(req *pb.SearchProductsRequest, headers *map[string]string) (*pb.SearchProductsResponse, error) {
    time.Sleep(extraLatency)

    var ps []*pb.Product
    for _, product := range p.parseCatalog() {
        if strings.Contains(strings.ToLower(product.Name), strings.ToLower(req.Query)) ||
            strings.Contains(strings.ToLower(product.Description), strings.ToLower(req.Query)) {
            ps = append(ps, product)
        }
    }

    return &pb.SearchProductsResponse{Results: ps}, nil
}

func (p *productCatalog) parseCatalog() []*pb.Product {
    if reloadCatalog || len(p.catalog.Products) == 0 {
        err := loadCatalog(&p.catalog)
        if err != nil {
            return []*pb.Product{}
        }
    }

    return p.catalog.Products
}
