# Service Architecture

This document explains how a gRPC service, a gRPC client, and a web service work in this architecture.

## gRPC service

In short, gRPC servers are replaced with the new architecture, which relies on protobuf marshalling.

### Constants and Variables

- `running_in_lambda`: This boolean variable states whether the service is running in AWS Lambda or not. If the value is
  true, the service will expect an object of the type `RequestData` and will return an object of the
  type `ResponseData`. If the value is false, the service will start an HTTP server and work with HTTP requests and
  responses. The value of this variable is true if and only if an environment variable with the name of `RUN_LAMBDA` and
  the value of `1` is present.
- `default_port`: This constant is the default port to be used when starting an HTTP server. It can be overwritten with
  an environment variable named `PORT`.
- RPC constants: These are string constants that have a name of the form `<rpc>_RPC`, with the `<rpc>` section being the
  name of an RPC function in the proto service definition. These constants are used to indicate which RPC function the
  client wanted to call.

### Classes and Structs

- `RequestData`: This class represents requests the Lambda function may receive. An object of this class can be
  initialized by deserializing
  a [JSON request payload string](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-integrations-lambda.html)
  (see version 2).
- `ResponseData`: This class represents responses the Lambda function may return. The objects of this class are
  serialized to JSON to be returned to the invoker. See the above link for the response format.

### Functions

- `call_rpc`: receives a proto message object and a `RequestData` object and invokes the appropriate RPC function based
  on the `rpc-name` header. The `RequestData` object is only used to retrieve the headers. It's assumed that
  the `rpc-name` header has a valid value since the `determine_message_type` function handles invalid values, and it's
  always invoked first. This function returns the exact output of the RPC function it invokes, which is a proto message
  and an RPC error.
- `determine_message_type`: determines the input message type based on the `rpc-name` header. The output of this
  function is used by the `decode_request` function to correctly unmarshall the proto message. It returns null if the
  header has an invalid value.
- `decode_request`: receives a `RequestData` object, reads the request body, optionally decodes it from base 64, and
  unmarshalls the byte data to a proto message based on the `rpc-name` header. If successful, it returns the message
  object. If not, it either returns a `ResponseData` object or raises an error. If the failure was the client's fault (a
  bad request) then a `ResponseData` object is returned, containing the gRPC status code and error message. Otherwise,
  an error is raised.
- `encode_response`: receives either a proto message object or an RPC error and returns a `ResponseData` object. Exactly
  one of the inputs has to be null. If the RPC error input is not null (meaning an error has occurred while running the
  RPC function) the output will contain a non-OK gRPC status and the error message. If not, the message object will be
  marshalled and a `ResponseData` object will be constructed based on the `running_in_lambda` variable and returned.
- `generate_error_eesponse`: receives a non-OK gRPC status code and error message string, and constructs and returns
  a `ResponseData` object containing the error data. It's used by the decode and encode functions.
- `main`: checks the `running_in_lambda` and invokes appropriate function.
- `run_lambda`: the main Lambda handler. This function is invoked only if the service is running in a Lambda function.
  This function receives a `RequestData`, decodes the request, returns the `ResponseData` object in case of a bad
  request, or calls the appropriate RPC function, encodes the response or the RPC error, and finally returns it.
- `run_http_server`: the main HTTP handler. This function is invoked only if the service is not running in a Lambda
  function. It defines a handler function which has the exact same functionally as the Lambda handler, and starts an
  HTTP server to serve on the specified address and port.
