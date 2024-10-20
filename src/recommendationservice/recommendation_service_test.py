import unittest

from common import GrpcError
from genproto import demo_pb2 as pb
from logger import getJSONLogger
from recommendation_service_stub import list_recommendations

logger = getJSONLogger('recommendationservice-test')


class TestRecommendationService(unittest.TestCase):
    def test_list_recommendations(self):
        try:
            response = list_recommendations(
                pb.ListRecommendationsRequest(user_id="test", product_ids=["test"]),
            )
            print(response.product_ids)
        except GrpcError as err:
            self.fail(f"RPC call failed with error: {err}")
        print(f"Test complete")


if __name__ == "__main__":
    unittest.main()
