const {HelloRequest, HelloResponse, ByeRequest, ByeResponse} = require('./genproto');

function sayHello(helloRequest, headers) {
    console.error(`Received: ${helloRequest.getName()}`);

    const helloResp = new HelloResponse();
    helloResp.setText(`Hello ${helloRequest.getName()}`);

    return helloResp;
}

function sayBye(byeRequest, headers) {
    console.error(`Received: ${byeRequest.getName()}`);

    const byeResp = new ByeResponse();
    byeResp.setText(`Bye ${byeRequest.getName()}`);

    return byeResp;
}

module.exports = {sayHello, sayBye};
