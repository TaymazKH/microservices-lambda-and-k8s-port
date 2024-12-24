# Development Guide

This is a guide to port existing services to this new architecture or develop new ones.

See the `/examples` directory. It contains some basic services that
demonstrate [our architecture](./service-architecture.md) in its simplest form. These services can be used as templates
for developing new services.

## Examples

There provided example services are structured in this way:

- `grpc_service`: This directory contains examples for gRPC services and clients in different languages. All the
  services are simple greeter services, as seen by their shared proto service definition in the `protos` directory.
    - `grpc_service/<lang>/server`: This directory contains a simple gRPC service in the specified language. You may
      deploy this service as instructed in the [deployment guide](./deployment-guide.md).
    - `grpc_service/<lang>/client`: This directory contains a simple gRPC client in the specified language. Note that
      these clients are not made to be deployed in AWS Lambda or Docker, rather they are for testing purposes as
      instructed in the [testing guide](./testing.md).
- `web_service`: This directory contains examples for web services in different languages. All the services are simple
  websites with only few pages.

## Development

To start developing your own services you may follow these instructions:

### gRPC Services

Follow these simple instructions to develop a gRPC service with minimal effort:

1. Start by copying the contents of the `grpc_service/<lang>/server` directory.
2. Modify the `protodir`, `protoname`, and proto file name in the `genproto.sh` script.
3. Rename the `greeter_service` file, and have your RPC functions in it. Implement your service logic in the same file
   or other files.
4. Edit the server file with minimal changes:
    1. Edit the constants: add your own RPC names.
    2. Edit the `call_rpc` function: for every possible case of RPC name (except the last one), call the appropriate RPC
       function.
    3. Edit the `determine_message_type` function: for every possible case of RPC name, return the appropriate request
       message class/object, or null, for an invalid value.
    4. \[Optional\] Add any other variable, constant, or piece of code you may need.
5. Edit the files for your desired deployment method:
    - Lambda: update the `deployment.sh` script to contain your desired files in the deployment package.
    - Docker: edit the service name and entrypoint in the `Dockerfile`. Optionally, edit the `.dockerignore` file to
      ignore any file you don't want in your container/image.

Your service is now ready to be deployed. You may choose to further work on your service of course.

### gRPC Clients

Follow these simple instructions to develop a gRPC client with minimal effort:

1. Start by copying the contents of the `grpc_service/<lang>/client` directory, except the `main` file.
2. Modify the `genproto.sh` script.
3. Edit the stub file(s):
    1. Rename the `greeter_service_stub` file.
    2. Rename/add RPC functions. Be careful to match the types if you're working on a typed language.
    3. Rename the address and timeout variables and environment variables.
4. Edit the RPC name constants. Depending on the language, these rae located either in the stubs or the client file.
   Correct the imports after this step.
5. Edit the `determine_message_type` function in the client file. For every possible case of RPC name, return the
   appropriate response message class/object.

Your client is now ready. You may choose to further work on your service of course.

### Web Services

Follow these simple instructions to develop a web service with minimal effort:

1. Start by copying the contents of the `web_service/<lang>` directory.
2. Modify the `genproto.sh` script.
3. Implement your service logic the way you see fit. Design the HTML and CSS files.
4. Edit the server file with minimal changes:
    1. Edit the `init` function: attach URLs to handlers.
    2. \[Optional\] If your code needs to work with multi-value headers that violate the comma-separated rule, you need
       to implement a parser for that. For this purpose, edit the `non_split_headers` variable and part of
       the `reconstruct_http_request` function that handles header parsing.
5. Edit the deployment script or the docker-related files for your desired deployment method.

Your service is now ready to be deployed. You may choose to further work on your service of course.
