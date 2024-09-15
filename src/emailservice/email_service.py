import grpc
from google.api_core.exceptions import GoogleAPICallError
from jinja2 import Environment, FileSystemLoader, select_autoescape, TemplateError

from genproto import demo_pb2
from genproto import demo_pb2_grpc
from logger import getJSONLogger

logger = getJSONLogger('emailservice-server')

# Loads confirmation email template from file
env = Environment(
    loader=FileSystemLoader('templates'),
    autoescape=select_autoescape(['html', 'xml'])
)
template = env.get_template('confirmation.html')


# ATTENTION: The main service logic is not implemented in the main benchmark.
# Therefore, we will use the dummy service.


class EmailService(demo_pb2_grpc.EmailServiceServicer):
    def __init__(self):
        raise Exception('cloud mail client not implemented')
        super().__init__()

    @staticmethod
    def send_email(client, email_address, content):
        response = client.send_message(
            sender=client.sender_path(project_id, region, sender_id),
            envelope_from_authority='',
            header_from_authority='',
            envelope_from_address=from_address,
            simple_message={
                "from": {
                    "address_spec": from_address,
                },
                "to": [{
                    "address_spec": email_address
                }],
                "subject": "Your Confirmation Email",
                "html_body": content
            }
        )
        logger.info("Message sent: {}".format(response.rfc822_message_id))

    def SendOrderConfirmation(self, request, context):
        email = request.email
        order = request.order

        try:
            confirmation = template.render(order=order)
        except TemplateError as err:
            context.set_details("An error occurred when preparing the confirmation mail.")
            logger.error(err.message)
            context.set_code(grpc.StatusCode.INTERNAL)
            return demo_pb2.Empty()

        try:
            EmailService.send_email(self.client, email, confirmation)
        except GoogleAPICallError as err:
            context.set_details("An error occurred when sending the email.")
            print(err.message)
            context.set_code(grpc.StatusCode.INTERNAL)
            return demo_pb2.Empty()

        return demo_pb2.Empty()
