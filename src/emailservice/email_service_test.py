import unittest

from client import send_order_confirmation
from common import GrpcError
from genproto import demo_pb2 as pb
from logger import getJSONLogger

logger = getJSONLogger('emailservice-test')


class TestEmailService(unittest.TestCase):
    def test_send_confirmation_email(self):
        try:
            empty = send_order_confirmation(
                send_order_confirmation_request=pb.SendOrderConfirmationRequest(email="example@gmail.com", order=None),
                addr="",
                timeout=5
            )
        except GrpcError as err:
            self.fail(f"RPC call failed with error: {err}")
        print(f"Test complete")


if __name__ == "__main__":
    unittest.main()
