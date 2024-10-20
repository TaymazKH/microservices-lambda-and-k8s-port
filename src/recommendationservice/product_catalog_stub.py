import os

from client import marshal_request, unmarshal_response, send_request, DEFAULT_TIMEOUT, \
    PRODUCT_CATALOG_SERVICE, LIST_PRODUCTS_RPC, GET_PRODUCT_RPC, SEARCH_PRODUCTS_RPC
from logger import getJSONLogger

logger = getJSONLogger('recommendationservice-productcatalog-stub')


def list_products(request, headers=None):
    bin_req = marshal_request(request)
    resp_body, headers = send_request(addr, PRODUCT_CATALOG_SERVICE, LIST_PRODUCTS_RPC, bin_req, headers, timeout)
    return unmarshal_response(resp_body, headers, LIST_PRODUCTS_RPC)


def get_product(request, headers=None):
    bin_req = marshal_request(request)
    resp_body, headers = send_request(addr, PRODUCT_CATALOG_SERVICE, GET_PRODUCT_RPC, bin_req, headers, timeout)
    return unmarshal_response(resp_body, headers, GET_PRODUCT_RPC)


def search_products(request, headers=None):
    bin_req = marshal_request(request)
    resp_body, headers = send_request(addr, PRODUCT_CATALOG_SERVICE, SEARCH_PRODUCTS_RPC, bin_req, headers, timeout)
    return unmarshal_response(resp_body, headers, SEARCH_PRODUCTS_RPC)


addr = os.getenv("PRODUCT_CATALOG_SERVICE_ADDR")
if addr is None:
    logger.error("PRODUCT_CATALOG_SERVICE_ADDR environment variable not set")
    raise EnvironmentError("PRODUCT_CATALOG_SERVICE_ADDR environment variable not set")

t = os.getenv("PRODUCT_CATALOG_SERVICE_TIMEOUT")
if t is None:
    timeout = DEFAULT_TIMEOUT
else:
    try:
        timeout = int(t)
        if timeout <= 0:
            timeout = DEFAULT_TIMEOUT
    except ValueError:
        timeout = DEFAULT_TIMEOUT
