from genproto import demo_pb2
from logger import getJSONLogger

logger = getJSONLogger('emailservice-dummylogic')


# ATTENTION: The main service logic is not implemented in the main benchmark.
# Therefore, we will use the dummy service.


class DummyEmailService:
    @staticmethod
    def SendOrderConfirmation(request, headers):
        logger.info('A request to send order confirmation email to {} has been received.'.format(request.email))
        return demo_pb2.Empty()
