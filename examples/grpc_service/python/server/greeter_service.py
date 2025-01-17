from genproto import greeter_pb2 as pb
from logger import get_simple_logger

logger = get_simple_logger('greeterservice-logic')


def say_hello(request, headers):
    logger.info(f"Received: {request.name}")
    return pb.HelloResponse(text=f"Hello {request.name}")


def say_bye(request, headers):
    logger.info(f"Received: {request.name}")
    return pb.HelloResponse(text=f"Bye {request.name}")
