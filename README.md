# Microservices Lambda and K8S Port

> How can we port gRPC services to a serverless environment, such as AWS Lambda?

This project aims to answer to this question.

Simply uploading the codebase to a Lambda function won't work, as there are some challenges on the way:

- AWS Lambda works only on HTTP/1 while gRPC uses HTTP/2 and its exclusive features.
- Lambda's serverless nature can’t handle the server-based nature of gRPC. The gRPC servers expect to receive HTTP
  requests, while in Lambda requests are handed in as "events" after invocation.

We have created a method of porting existing services with minimal modifications to the code base, and developing new
ones with at-most simplicity. Our method utilizes a custom service architecture which is easily attached to gRPC
services, alongside Protocol Buffer marshalling to communicate data. The services in the new architecture have the
ability to be executed in multiple manners, either in Lambda, a Docker container, or simply on your machine.

Further, we have demonstrated our creation by porting a web service which consists of multiple gRPC microservices. We
used [GoogleCloudPlatform/microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo) for this
purpose.

In short, our contributions include:

- A method of porting gRPC servers and clients to be able to be executed both in Lambda functions and Docker containers.
- A method of porting HTTP servers to be able to be executed in Lambda functions as well as Docker containers.
- Examples demonstrating simple servers and clients ported to the architecture, alongside detailed documents.
- Ported the [GoogleCloudPlatform/microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo)
  repository as our main benchmark, to demonstrate how a whole system would behave.

## Next Steps

Detailed documents can be found under the [`/docs`](./docs) directory.

- To deploy the service, see [deployment guide](./docs/deployment-guide.md).
- To test a singular microservice automatically, see [testing guide](./docs/testing.md).
- To quickly develop your own services, see [development guide](./docs/development-guide.md).
- To learn more about our service architecture, see [service architecture](./docs/service-architecture.md).
