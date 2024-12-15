# Testing Guide

You can find the tests in the [`/tests`](../tests) directory.

**Note 1:** The tests are meant to check the correct deployment and connectivity of the services. They do not have full
coverage nor check the correctness of the service logics. It's assumed that the service logics work correctly.

**Note 2:** The frontend service has no tests as it's intended to be tested with a browser.

## How to test a service

1. Deploy the intended service and its dependency services. See [deployment guide](./deployment-guide.md).
2. Set the service's address as an environment variables.
3. Go to the specific service's test directory and run the test.
    - Golang: `go test -v -count=1 ./client`.
    - JS: `node <test>.js`.
    - Python: `python <test>.py`.
