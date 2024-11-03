const {HelloResponse, ByeResponse} = require('./genproto/hello_pb');

function sayHello(helloRequest, headers) {
    console.log(`Received: ${helloRequest.getName()}`);

    const helloResp = new HelloResponse();
    helloResp.setText(`Hello ${helloRequest.getName()}`);

    return helloResp;
}

function sayBye(byeRequest, headers) {
    console.log(`Received: ${byeRequest.getName()}`);

    const byeResp = new ByeResponse();
    byeResp.setText(`Bye ${byeRequest.getName()}`);

    return byeResp;
}

module.exports = {sayHello, sayBye};
