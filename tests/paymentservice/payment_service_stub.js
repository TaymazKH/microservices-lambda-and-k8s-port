const {
    marshalRequest,
    unmarshalResponse,
    sendRequest,
    DEFAULT_TIMEOUT,
    PAYMENT_SERVICE,
    CHARGE_RPC
} = require('./client');

async function charge(request, headers = {}) {
    const binReq = marshalRequest(request);
    const [respBody, respHeaders] = await sendRequest(paymentServiceAddr, PAYMENT_SERVICE, CHARGE_RPC, binReq, headers, paymentServiceTimeout);
    return unmarshalResponse(respBody, respHeaders, CHARGE_RPC);
}

const paymentServiceAddr = process.env.PAYMENT_SERVICE_ADDR;
if (!paymentServiceAddr) {
    console.error("PAYMENT_SERVICE_ADDR environment variable not set");
    throw new Error("PAYMENT_SERVICE_ADDR environment variable not set");
}

let t = process.env.PAYMENT_SERVICE_TIMEOUT;
if (!t) {
    t = DEFAULT_TIMEOUT;
} else {
    t = parseInt(t, 10);
    if (isNaN(t) || t <= 0) {
        t = DEFAULT_TIMEOUT;
    }
}
const paymentServiceTimeout = t

module.exports = {charge};
