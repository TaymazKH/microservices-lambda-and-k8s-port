# Service Architecture

This document explains how a gRPC service, a gRPC client, and a web service work in this architecture.

## gRPC Service

In short, gRPC servers are replaced with the new architecture, which relies on protobuf marshalling.

### Constants and Variables

- `running_in_lambda`: This boolean variable states whether the service is running in AWS Lambda or not. If the value is
  true, the service will expect an object of the type `RequestData` and will return an object of the
  type `ResponseData`. If the value is false, the service will start an HTTP server and work with HTTP requests and
  responses. The value of this variable is true if and only if an environment variable with the name of `RUN_LAMBDA` and
  the value of `1` is present.
- `default_port`: This constant is the default port to be used when starting an HTTP server. It can be overwritten with
  an environment variable named `PORT`.
- RPC name constants: These are string constants that have a name of the form `<RPC_NAME>_RPC`, with the `<RPC_NAME>`
  section being the name of an RPC function in the proto service definition. These constants are used to indicate which
  RPC function the client wanted to call.

### Classes and Structs

- `RequestData`: This class represents requests the Lambda function may receive. An object of this class can be
  initialized by deserializing
  a [JSON request payload string](https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-integrations-lambda.html)
  (see version 2).
- `ResponseData`: This class represents responses the Lambda function may return. The objects of this class are
  serialized to JSON to be returned to the invoker. See the above link for the response format.

### Server Functions

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
- `main`: checks the `running_in_lambda` variable and invokes the appropriate function.
- `run_lambda`: the main Lambda handler. This function is invoked only if the service is running in a Lambda function.
  This function receives a `RequestData`, decodes the request, returns the `ResponseData` object in case of a bad
  request, or calls the appropriate RPC function, encodes the response or the RPC error, and finally returns it.
- `run_http_server`: the main HTTP handler. This function is invoked only if the service is not running in a Lambda
  function. It defines a handler function which has the exact same functionally as the Lambda handler, and starts an
  HTTP server to serve on the specified address and port.

### RPC Functions

These functions are located in a file with the same name as the service. These are the main RPC handlers and each
correspond to an RPC in the proto service definition. These functions must have this structure:

- Inputs:
    1. A proto request message. All the RPC functions must have only one message input. If an RPC requires more than one
       argument, a new message type must be defined, containing all the arguments.
    2. A header map. The headers are used for communicating context between the client and the server, such as the
       timeout. Services may define their own structure for communicating context, such as defining custom headers.
- Outputs:
    1. A proto response message. All the RPC functions must have only one message output. If an RPC requires more than
       one argument, a new message type must be defined.
- Errors:
    1. A gRPC error. A function may raise a gRPC error containing a gRPC status code and an error message. Services may
       use the error message to communicate external data and traces, for example, by using a JSON string as the
       message. For each language, the following methods must be used to instantiate an error:
        - Golang: The `Error` and `Errorf` functions, imported from `google.golang.org/grpc/status`.
        - JS: Simply by instantiating an object with two fields named `code` and `message` (like
          this: `throw {code: status.NOT_FOUND, message: "Not Found"};`). You may import and use the `status` package
          from `@grpc/grpc-js`.
        - Python: By using the `GrpcError` class in the `common.py` file that is provided in every Python service.

## gRPC Client

### Constants and Variables

- `default_timeout`: the default timeout interval when sending a request. It can be overwritten by an environment
  variable.

### Functions

- `determine_message_type`: determines the output message type based on the `rpc_name` string value. The output of this
  function is used by the `unmarshal_response` function to correctly unmarshall the proto message. It's assumed that
  the RPC name is a valid value.
- `send_request`: receives the address, service information, binary request data, headers, and some other data to send
  an HTTP request and return the response body and headers. It sets the `rpc-name` and `content-type` headers and sends
  the binary data as a POST request to `<addr>/<service_name>` with the given timeout. It reads the response body and
  headers and returns them. An error is raised in case of any failure or non 200 status code.
- `marshal_request`: simply marshalls the proto message object and returns the binary data.
- `unmarshal_response`: receives the outputs of the `send_request` function and the RPC name, and either returns the
  unmarshalled response or raises an error, which may be an RPC error that's returned by the server or an error that has
  occurred in this function. An error is raised if the `grpc-status` header is missing or has an invalid value, or the
  response body can't be unmarshalled. If the gRPC status code indicates success the binary data will be unmarshalled to
  a proto message based on the output of the `determine_message_type` function. If not, the RPC error will be
  reconstructed with the code and the error message.

### Stubs

A stub represents a service's client. A service has as many stubs as the services it depends on. A stub has these parts:

- Constants and Variables:
    - `<service_name>_addr`: is the address of the service. An environment variable with the name
      of `<SERVICE_NAME>_ADDR` must be present to populate this variable.
    - `<service_name>_timeout`: is the timeout interval used when calling this service's RPC. The default timeout will
      be used if the `<SERVICE_NAME>_TIMEOUT` environment variable isn't present or has an invalid value.
    - `<service_name>_service`: is a constant string that stores the name of the service. It's used as the part for
      sending HTTP requests.
    - RPC name constants: same as the constants with the same names in gRPC services.
- Functions: each function in a stub (except the `init` function in Golang stubs) represents an RPC. They receive a
  proto message (same as the RPC input in the service definition) and headers (as a way to communicate context), and
  either return a proto message or raise an error. They call the `marshal_request`, `send_request`,
  and `unmarshal_response` functions in order.

## Web Service

In our architecture, web services are created directly from normal HTTP servers by attaching the Lambda integration to
the existing code. In Lambda, HTTP requests are reconstructed from Lambda requests, then the existing HTTP handlers are
invoked, and finally the HTTP response is captured and converted to a Lambda Response. This was to make porting existing
services much easier.

### Constants and Variables

- `running_in_lambda`: same as the variable with the same name in gRPC services.
- `base_url`: is the base URL of the service. It's optionally populated with the `BASE_URL` environment variable. It
  must be either empty or be of the form of a valid path, beginning with a slash.
- `http_handler`: is the main handler function of whole web service. It must map all the valid URLs to a handler
  function. It can be a nested handler for middleware purposes, such as logging and handling sessions.
- `default_port`: same as the constant with the same name in gRPC services.
- `non_split_headers`: is a set of header names that violate the comma-separation rule in multi-value headers. These
  headers can't be simply split based on a comma and need more precise handling.

### Classes and Structs

- `RequestData` and `ResponseData`: same as the classes with the same names in gRPC services. However, these have more
  fields as a web service may need more data to serve.

### Server Functions

- `reconstruct_http_request`: receives a `RequestData` object and returns an HTTP request object. It initializes an HTTP
  request object and populates its fields based on the fields of the `RequestData` object.
- `convert_to_response_data`: receives a `ResponseData` object and returns an HTTP response object. It initializes an
  HTT response object and populates its fields based on the fields of the `ResponseData` object.
- `main`: checks the `running_in_lambda` variable and invokes the appropriate function.
- `run_lambda`: the main Lambda handler. This function is invoked only if the service is running in a Lambda function.
  This function receives a `RequestData`, reconstructs the HTTP request, initializes a response recorder, invokes the
  HTTP handler with the response recorder as the writer, retrieves the HTTP response from the recorder, converts it to
  a `ResponseData` object, and finally returns it.
- `run_http_server`: the main HTTP handler. This function is invoked only if the service is not running in a Lambda
  function. It starts an HTTP server with the `http_handler` variable as its handler, to serve on the specified address
  and port.

### Handler Functions

These functions are the main HTTP handlers and located in a file named `handlers`. They have the same structure as
typical HTTP handlers in a server (in case of Golang, this means they take a writer and a request as arguments). 
