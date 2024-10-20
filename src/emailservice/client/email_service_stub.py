import os

from .client import marshal_request, unmarshal_response, send_request, DEFAULT_TIMEOUT, \
    EMAIL_SERVICE, SEND_ORDER_CONFIRMATION_RPC
from ..logger import getJSONLogger

logger = getJSONLogger('emailservice-stub')


def send_order_confirmation(request, headers=None):
    bin_req = marshal_request(request)
    resp_body, headers = send_request(greeter_addr, EMAIL_SERVICE, SEND_ORDER_CONFIRMATION_RPC, bin_req, headers,
                                      greeter_timeout)
    return unmarshal_response(resp_body, headers, SEND_ORDER_CONFIRMATION_RPC)


greeter_addr = os.getenv("EMAIL_SERVICE_ADDR")
if greeter_addr is None:
    logger.error("EMAIL_SERVICE_ADDR environment variable not set")
    raise EnvironmentError("EMAIL_SERVICE_ADDR environment variable not set")

t = os.getenv("EMAIL_SERVICE_TIMEOUT")
if t is None:
    greeter_timeout = DEFAULT_TIMEOUT
else:
    try:
        greeter_timeout = int(t)
        if greeter_timeout <= 0:
            greeter_timeout = DEFAULT_TIMEOUT
    except ValueError:
        greeter_timeout = DEFAULT_TIMEOUT
