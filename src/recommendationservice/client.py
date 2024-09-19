import grpc
import requests
from google.protobuf.message import DecodeError

from common import GrpcError
from genproto import demo_pb2 as pb

RECOMMENDATION_SERVICE = "recommendation-service"
LIST_RECOMMENDATIONS_RPC = "list-recommendations"


def list_recommendations(list_recommendations_request, addr, timeout):
    bin_req = marshal_request(list_recommendations_request)

    path = f"/{RECOMMENDATION_SERVICE}/{LIST_RECOMMENDATIONS_RPC}"
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
        msg = pb.ListRecommendationsResponse()

        try:
            msg.ParseFromString(resp_body)
            return msg
        except DecodeError as e:
            raise ValueError(f"Failed to parse response body: {e}")
    else:
        raise GrpcError(code=grpc_code, message=resp_body)
