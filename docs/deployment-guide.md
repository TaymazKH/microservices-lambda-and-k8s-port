# Deployment Guide

You can deploy each service either on AWS Lambda, in a docker container, or on your local machine.

## AWS Lambda

To deploy a service in a Lambda function follow these steps:

1. Make sure you have the means to compile proto files.
    - Golang & JS: Install protocol buffer compiler. Follow this [tutorial](https://grpc.io/docs/protoc-installation/).
    - Python: Install `grpcio_tools` using pip or any other package manager.
2. Run the `genproto.sh` script to compile to proto file.
3. Run the `deployment.sh` script to create the deployment package, named `deployment.zip`.
4. Create a Lambda function. Make sure to use the correct runtime and give it the needed permissions. Use the `amd64`
   architecture. Do the basic configurations such as the timeout value.
5. Upload the deployment package.
6. Set the handler function's name.
    - Golang: Not needed.
    - JS: `server/runLambda`.
    - Python: `server/run_lambda`.
7. Add a new environment variable named `RUN_LAMBDA` with the value of `1`. Also, add any other variables that may be
   needed by the service. Refer to the readme file in a specific service's directory.
8. Set up the means to invoke this function. This may be a function URL or AWS API gateway integration.

## Docker

To deploy a service in a docker container follow these steps:

1. Do the first two steps of the [AWS Lambda](#aws-lambda) section.
2. A service that depends on other services uses environment variables to store their addresses. You may either...
    - ... replace the addresses in the docker file.
    - ... override the addresses when running a container.
    - ... create a docker network and deploy the services in containers, using the names as seen in the docker files, in
      the network.
3. Build the image based on the docker file.
4. Run a container based on that image. Optionally, override any of the environment variables or set up a network.

## Local

To deploy a service locally follow these steps:

1. Do the first two steps of the [AWS Lambda](#aws-lambda) section.
2. Install the dependencies.
    - Golang: `go mod download`.
    - JS: `npm install`.
    - Python: `pip install -r requirements.txt`.
3. Set the required environment variables (see the readme file in a specific service's directory).
4. Run the service.
    - Golang: `go run .`. You may also build the service and then run it.
    - JS: `node server.js`.
    - Python: `python server.py`.
