import os

import grpc
import requests
from google.protobuf.message import DecodeError

from common import GrpcError
from genproto import demo_pb2 as pb
from logger import getJSONLogger

logger = getJSONLogger('recommendationservice-server')

PRODUCT_CATALOG_SERVICE = "product-catalog-service"
LIST_PRODUCTS_RPC = "list-products"
GET_PRODUCT_RPC = "get-product"
SEARCH_PRODUCTS_RPC = "search-products"

addr = os.getenv('PRODUCT_CATALOG_SERVICE_ADDR')
if addr is None:
    raise Exception('PRODUCT_CATALOG_SERVICE_ADDR environment variable not set')
logger.info("product catalog address: " + addr)

timeout = int(os.getenv('PRODUCT_CATALOG_SERVICE_TIMEOUT', '5'))


def list_products(empty):
    bin_req = marshal_request(empty)

    path = f"/{PRODUCT_CATALOG_SERVICE}/{LIST_PRODUCTS_RPC}"
    resp_body, headers = send_request(addr, path, bin_req, timeout)

    return unmarshal_response(resp_body, headers, path)


def get_product(get_product_request):
    bin_req = marshal_request(get_product_request)

    path = f"/{PRODUCT_CATALOG_SERVICE}/{GET_PRODUCT_RPC}"
    resp_body, headers = send_request(addr, path, bin_req, timeout)

    return unmarshal_response(resp_body, headers, path)


def search_products(search_products_request):
    bin_req = marshal_request(search_products_request)

    path = f"/{PRODUCT_CATALOG_SERVICE}/{SEARCH_PRODUCTS_RPC}"
    resp_body, headers = send_request(addr, path, bin_req, timeout)

    return unmarshal_response(resp_body, headers, path)


def send_request(addr, path, bin_req, timeout):
    try:
        response = requests.post(f"{addr}{path}",
                                 data=bin_req,
                                 headers={"Content-Type": "application/octet-stream"},
                                 timeout=timeout)
        response.raise_for_status()
        return response.content, response.headers
    except requests.RequestException as e:
        raise RuntimeError(f"Failed to send request: {e}")


def marshal_request(msg):
    return msg.SerializeToString()


def unmarshal_response(resp_body, headers, path):
    grpc_code = headers.get("Grpc-Code")
    if grpc_code is None:
        raise ValueError("Missing Grpc-Code header")

    grpc_code = int(grpc_code)

    if grpc_code == grpc.StatusCode.OK.value[0]:
        if path == f"/{PRODUCT_CATALOG_SERVICE}/{LIST_PRODUCTS_RPC}":
            msg = pb.ListProductsResponse()
        elif path == f"/{PRODUCT_CATALOG_SERVICE}/{GET_PRODUCT_RPC}":
            msg = pb.Product()
        else:
            msg = pb.SearchProductsResponse()

        try:
            msg.ParseFromString(resp_body)
            return msg
        except DecodeError as e:
            raise ValueError(f"Failed to parse response body: {e}")
    else:
        raise GrpcError(code=grpc_code, message=resp_body)
