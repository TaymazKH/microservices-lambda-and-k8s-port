import random

import product_catalog_stub
from genproto import demo_pb2 as pb
from logger import getJSONLogger

logger = getJSONLogger('recommendationservice-logic')


class RecommendationService:
    @staticmethod
    def ListRecommendations(request, headers):
        max_responses = 5

        # fetch list of products from product catalog stub
        cat_response = product_catalog_stub.list_products(pb.Empty())
        product_ids = [x.id for x in cat_response.products]
        filtered_products = list(set(product_ids) - set(request.product_ids))
        num_products = len(filtered_products)
        num_return = min(max_responses, num_products)

        # sample list of indices to return
        indices = random.sample(range(num_products), num_return)

        # fetch product ids from indices
        prod_list = [filtered_products[i] for i in indices]
        logger.info("[Recv ListRecommendations] product_ids={}".format(prod_list))

        # build and return response
        response = pb.ListRecommendationsResponse()
        response.product_ids.extend(prod_list)
        return response
