import os

from client import marshal_request, unmarshal_response, send_request, DEFAULT_TIMEOUT, \
    GREETER_SERVICE, SAY_HELLO_RPC, SAY_BYE_RPC
from logger import get_simple_logger

logger = get_simple_logger('greeterservice-stub')


def say_hello(request, headers=None):
    bin_req = marshal_request(request)
    resp_body, headers = send_request(greeter_addr, GREETER_SERVICE, SAY_HELLO_RPC, bin_req, headers, greeter_timeout)
    return unmarshal_response(resp_body, headers, SAY_HELLO_RPC)


def say_bye(request, headers=None):
    bin_req = marshal_request(request)
    resp_body, headers = send_request(greeter_addr, GREETER_SERVICE, SAY_BYE_RPC, bin_req, headers, greeter_timeout)
    return unmarshal_response(resp_body, headers, SAY_BYE_RPC)


greeter_addr = os.getenv("GREETER_SERVICE_ADDR")
if greeter_addr is None:
    logger.error("GREETER_SERVICE_ADDR environment variable not set")
    raise EnvironmentError("GREETER_SERVICE_ADDR environment variable not set")

t = os.getenv("GREETER_SERVICE_TIMEOUT")
if t is None:
    greeter_timeout = DEFAULT_TIMEOUT
else:
    try:
        greeter_timeout = int(t)
        if greeter_timeout <= 0:
            greeter_timeout = DEFAULT_TIMEOUT
    except ValueError:
        greeter_timeout = DEFAULT_TIMEOUT
