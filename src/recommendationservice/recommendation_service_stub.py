import os

from client import marshal_request, unmarshal_response, send_request, DEFAULT_TIMEOUT, \
    RECOMMENDATION_SERVICE, LIST_RECOMMENDATIONS_RPC
from logger import getJSONLogger

logger = getJSONLogger('recommendationservice-stub')


def list_recommendations(request, headers=None):
    bin_req = marshal_request(request)
    resp_body, headers = send_request(addr, RECOMMENDATION_SERVICE, LIST_RECOMMENDATIONS_RPC, bin_req, headers, timeout)
    return unmarshal_response(resp_body, headers, LIST_RECOMMENDATIONS_RPC)


addr = os.getenv("RECOMMENDATION_SERVICE_ADDR")
if addr is None:
    logger.error("RECOMMENDATION_SERVICE_ADDR environment variable not set")
    raise EnvironmentError("RECOMMENDATION_SERVICE_ADDR environment variable not set")

t = os.getenv("RECOMMENDATION_SERVICE_TIMEOUT")
if t is None:
    timeout = DEFAULT_TIMEOUT
else:
    try:
        timeout = int(t)
        if timeout <= 0:
            timeout = DEFAULT_TIMEOUT
    except ValueError:
        timeout = DEFAULT_TIMEOUT
