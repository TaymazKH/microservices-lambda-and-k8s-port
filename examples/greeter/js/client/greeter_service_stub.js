const {
    marshalRequest,
    unmarshalResponse,
    sendRequest,
    DEFAULT_TIMEOUT,
    GREETER_SERVICE,
    SAY_HELLO_RPC,
    SAY_BYE_RPC
} = require('./client');

async function sayHello(request, headers = {}) {
    const binReq = marshalRequest(request);
    const [respBody, respHeaders] = await sendRequest(greeterAddr, GREETER_SERVICE, SAY_HELLO_RPC, binReq, headers, greeterTimeout);
    return unmarshalResponse(respBody, respHeaders, SAY_HELLO_RPC);
}

async function sayBye(request, headers = {}) {
    const binReq = marshalRequest(request);
    const [respBody, respHeaders] = await sendRequest(greeterAddr, GREETER_SERVICE, SAY_BYE_RPC, binReq, headers, greeterTimeout);
    return unmarshalResponse(respBody, respHeaders, SAY_BYE_RPC);
}

const greeterAddr = process.env.GREETER_SERVICE_ADDR;
if (!greeterAddr) {
    console.error("GREETER_SERVICE_ADDR environment variable not set");
    throw new Error("GREETER_SERVICE_ADDR environment variable not set");
}

let t = process.env.GREETER_SERVICE_TIMEOUT;
if (!t) {
    t = DEFAULT_TIMEOUT;
} else {
    t = parseInt(t, 10);
    if (isNaN(t) || t <= 0) {
        t = DEFAULT_TIMEOUT;
    }
}
const greeterTimeout = t

module.exports = {sayHello, sayBye};
