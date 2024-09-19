import unittest

from client import list_recommendations
from common import GrpcError
from genproto import demo_pb2 as pb
from logger import getJSONLogger

logger = getJSONLogger('recommendationservice-server')


class TestRecommendationService(unittest.TestCase):
    def test_list_recommendations(self):
        try:
            response = list_recommendations(
                list_recommendations_request=pb.ListRecommendationsRequest(user_id="test", product_ids=["test"]),
                addr="",
                timeout=5
            )
            print(response.product_ids)
        except GrpcError as err:
            self.fail(f"RPC call failed with error: {err}")
        print(f"Test complete")


if __name__ == "__main__":
    unittest.main()
