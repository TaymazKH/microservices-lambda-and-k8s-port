from genproto import demo_pb2
from logger import getJSONLogger

logger = getJSONLogger('emailservice-server')


class DummyEmailService:
    @staticmethod
    def SendOrderConfirmation(request, context):
        logger.info('A request to send order confirmation email to {} has been received.'.format(request.email))
        return demo_pb2.Empty()
