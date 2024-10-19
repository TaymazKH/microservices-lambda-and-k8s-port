import base64
import os
from http.server import BaseHTTPRequestHandler, HTTPServer

import grpc
from google.protobuf.message import DecodeError, Message

from common import GrpcError
from genproto import hello_pb2 as pb
from greeter_service import say_hello, say_bye
from logger import get_simple_logger

logger = get_simple_logger('greeterservice-server')

RUNNING_IN_LAMBDA = os.getenv("RUN_LAMBDA") == "1"
DEFAULT_PORT = "8080"

SAY_HELLO_RPC = "say-hello"
SAY_BYE_RPC = "say-bye"


def call_rpc(msg, req_data):
    rpc_name = req_data.headers.get("rpc-name")
    if rpc_name == SAY_HELLO_RPC:
        return say_hello(msg, req_data.headers)
    else:
        return say_bye(msg, req_data.headers)


def determine_message_type(rpc_name):
    if rpc_name == SAY_HELLO_RPC:
        return pb.HelloRequest()
    elif rpc_name == SAY_BYE_RPC:
        return pb.ByeRequest()
    else:
        return None


class RequestData:
    def __init__(self, body, headers, isBase64Encoded, **kwargs):
        self.body = body
        self.headers = headers
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


def decode_request(req_data: RequestData) -> (Message, ResponseData):
    if req_data.isBase64Encoded:
        bin_req_body = base64.b64decode(req_data.body)
    else:
        bin_req_body = req_data.body

    rpc_name = req_data.headers.get('rpc-name')
    msg = determine_message_type(rpc_name)
    if not msg:
        return None, generate_error_response(grpc.StatusCode.UNIMPLEMENTED, f"Unknown RPC name: {rpc_name}")

    try:
        msg.ParseFromString(bin_req_body)
    except DecodeError as e:
        return None, generate_error_response(grpc.StatusCode.INVALID_ARGUMENT, str(e))

    return msg, None


def encode_response(msg: Message = None, rpc_error: GrpcError = None) -> ResponseData:
    if rpc_error:
        return generate_error_response(rpc_error.code, rpc_error.message)

    bin_resp_body = msg.SerializeToString()

    if RUNNING_IN_LAMBDA:
        return ResponseData(200, {
            "content-type": "application/octet-stream",
            "grpc-status": str(grpc.StatusCode.OK.value[0])
        }, base64.b64encode(bin_resp_body).decode('utf-8'), True)

    else:
        return ResponseData(200, {
            "content-type": "application/octet-stream",
            "grpc-status": str(grpc.StatusCode.OK.value[0])
        }, bin_resp_body, False)


def generate_error_response(code: grpc.StatusCode, message: str) -> ResponseData:
    return ResponseData(200, {
        "content-type": "text/plain",
        "grpc-status": str(code.value[0])
    }, message, False)


def run_lambda(event, context):
    logger.info("Handler started.")
    logger.info(f"Event data: {event}")

    req_data = RequestData(**event)
    req_msg, resp_data = decode_request(req_data)

    if resp_data is None:
        try:
            resp_msg = call_rpc(req_msg, req_data)
            resp_data = encode_response(msg=resp_msg)
        except GrpcError as err:
            resp_data = encode_response(rpc_error=err)

    json_response = resp_data.to_dict()

    logger.info(f"Response: {json_response}")
    logger.info("Handler finished.")
    return json_response


def run_http_server():
    class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):
        def do_POST(self):
            req_body = self.rfile.read(int(self.headers['Content-Length']))

            headers = {k.lower(): v for k, v in self.headers.items()}

            req_data = RequestData(req_body, headers, False)

            try:
                req_msg, resp_data = decode_request(req_data)

                if resp_data is None:
                    try:
                        resp_msg = call_rpc(req_msg, req_data)
                        resp_data = encode_response(msg=resp_msg)
                    except GrpcError as err:
                        resp_data = encode_response(rpc_error=err)

                self.send_response(resp_data.statusCode)
                for k, v in resp_data.headers.items():
                    self.send_header(k, v)
                self.end_headers()
                self.wfile.write(resp_data.body)

            except Exception as e:
                logger.error(f"Error handling request: {e}")
                self.send_response(500)
                self.end_headers()
                self.wfile.write(b"Internal Server Error")

    port = int(os.getenv("PORT", DEFAULT_PORT))
    logger.info(f"Port: {port}")

    httpd = HTTPServer(('', port), SimpleHTTPRequestHandler)
    httpd.serve_forever()


if __name__ == '__main__':
    if RUNNING_IN_LAMBDA:
        logger.warn("Conflict: RUN_LAMBDA=1 and __name__=__main__")
    else:
        run_http_server()
