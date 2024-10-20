import unittest

from .email_service_stub import send_order_confirmation
from ..common import GrpcError
from ..genproto import demo_pb2 as pb
from ..logger import getJSONLogger

logger = getJSONLogger('emailservice-test')


class TestEmailService(unittest.TestCase):
    def test_send_confirmation_email(self):
        try:
            empty = send_order_confirmation(
                pb.SendOrderConfirmationRequest(email="example@gmail.com", order=None),
            )
        except GrpcError as err:
            self.fail(f"RPC call failed with error: {err}")
        print(f"Test complete")


if __name__ == "__main__":
    unittest.main()
