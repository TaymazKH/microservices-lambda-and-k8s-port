import argparse
from genproto import hello_pb2 as pb
from greeter_service_stub import say_hello, say_bye

from logger import get_simple_logger

logger = get_simple_logger('main')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--name", type=str, default="world", help="Name to greet")
    args = parser.parse_args()
    name = args.name

    hello_req = pb.HelloRequest(name=name)

    try:
        hello_resp = say_hello(hello_req)
    except Exception as e:
        logger.error(f"Error calling SayHello RPC: {e}")
        return

    logger.info(f"Greeting: {hello_resp.text}")

    bye_req = pb.ByeRequest(name=name)

    try:
        bye_resp = say_bye(bye_req)
    except Exception as e:
        logger.error(f"Error calling SayBye RPC: {e}")
        return

    logger.info(f"Farewell: {bye_resp.text}")

    print(f"{hello_resp.text} - {bye_resp.text}")


if __name__ == "__main__":
    main()
