import grpc
import requests
from google.protobuf.message import EncodeError, DecodeError

from common import GrpcError
from genproto import hello_pb2 as pb

DEFAULT_TIMEOUT = 10

GREETER_SERVICE = "greeter"
SAY_HELLO_RPC = "say-hello"
SAY_BYE_RPC = "say-bye"


def determine_message_type(rpc_name):
    if rpc_name == SAY_HELLO_RPC:
        return pb.HelloResponse()
    else:
        return pb.ByeResponse()


def send_request(addr, service_name, rpc_name, bin_req, headers=None, timeout=DEFAULT_TIMEOUT):
    if headers is None:
        headers = {}
    headers['rpc-name'] = rpc_name
    headers['content-type'] = 'application/octet-stream'

    try:
        response = requests.post(f"{addr}/{service_name}", data=bin_req, headers=headers, timeout=timeout)
        response.raise_for_status()
        return response.content, response.headers
    except requests.RequestException as e:
        raise RuntimeError(f"Failed to send HTTP request: {e}")


def marshal_request(msg):
    try:
        return msg.SerializeToString()
    except EncodeError as e:
        raise RuntimeError(f"Failed to marshal request: {e}")


def unmarshal_response(resp_body, headers, rpc_name):
    grpc_status = headers.get("grpc-status")
    if grpc_status is None:
        raise KeyError("Missing grpc-status header")

    try:
        grpc_status = int(grpc_status)
    except ValueError:
        raise ValueError(f"Failed to parse grpc-status header: {grpc_status}")

    if grpc_status == grpc.StatusCode.OK.value[0]:
        msg = determine_message_type(rpc_name)

        try:
            msg.ParseFromString(resp_body)
            return msg
        except DecodeError as e:
            raise RuntimeError(f"Failed to unmarshal response: {e}")
    else:
        raise GrpcError(code=grpc_status, message=resp_body)
