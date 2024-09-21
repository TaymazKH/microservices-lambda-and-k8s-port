import base64

import grpc

from common import GrpcError
from dummy_email_service import DummyEmailService as Service
from genproto import demo_pb2 as pb
from logger import getJSONLogger

logger = getJSONLogger('emailservice-server')

EMAIL_SERVICE = "email-service"
SEND_ORDER_CONFIRMATION_RPC = "send-order-confirmation"


class RequestContext:
    def __init__(self, method, path, **kwargs):
        self.http = {"method": method, "path": path}


class RequestData:
    def __init__(self, body, headers, requestContext, isBase64Encoded, **kwargs):
        self.body = body
        self.headers = headers
        self.requestContext = requestContext
        self.isBase64Encoded = isBase64Encoded


class ResponseData:
    def __init__(self, statusCode, headers, body, isBase64Encoded):
        self.statusCode = statusCode
        self.headers = headers
        self.body = body
        self.isBase64Encoded = isBase64Encoded

    def to_dict(self):
        return {
            "statusCode": self.statusCode,
            "headers": self.headers,
            "body": self.body,
            "isBase64Encoded": self.isBase64Encoded,
        }


def handle_request(msg, req_data):
    return Service.SendOrderConfirmation(msg)


def decode_request(req_data):
    if req_data['isBase64Encoded']:
        bin_req_body = base64.b64decode(req_data['body'])
    else:
        bin_req_body = req_data['body'].encode()

    if req_data['requestContext']['http']['path'] == f"/{EMAIL_SERVICE}/{SEND_ORDER_CONFIRMATION_RPC}":
        msg = pb.SendOrderConfirmationRequest()
    else:
        raise ValueError(f"Unknown path: {req_data['requestContext']['http']['path']}")

    msg.ParseFromString(bin_req_body)
    return msg, RequestData(**req_data)


def encode_response(msg, rpc_error=None):
    if rpc_error is None:
        bin_resp_body = msg.SerializeToString()
        encoded_resp_body = base64.b64encode(bin_resp_body).decode()

        resp_data = ResponseData(
            statusCode=200,
            headers={"Content-Type": "application/octet-stream", "Grpc-Code": str(grpc.StatusCode.OK.value[0])},
            body=encoded_resp_body,
            isBase64Encoded=True
        )
    else:
        resp_data = ResponseData(
            statusCode=200,
            headers={"Content-Type": "text/plain", "Grpc-Code": str(rpc_error.int_code)},
            body=rpc_error.message,
            isBase64Encoded=False
        )

    return resp_data.to_dict()


def main(event, context):
    logger.info("Handler started.")
    logger.info(f"Event data: {event}")
    req_msg, req_data = decode_request(event)

    try:
        resp_msg = handle_request(req_msg, req_data)
        response = encode_response(resp_msg)
    except GrpcError as err:
        response = encode_response(None, rpc_error=err)

    logger.info(f"Response: {response}")
    logger.info("Handler finished.")
    return response
