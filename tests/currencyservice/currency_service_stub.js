const {
    marshalRequest,
    unmarshalResponse,
    sendRequest,
    DEFAULT_TIMEOUT,
    CURRENCY_SERVICE,
    GET_SUPPORTED_CURRENCIES_RPC,
    CONVERT_RPC
} = require('./client');

async function getSupportedCurrencies(request, headers = {}) {
    const binReq = marshalRequest(request);
    const [respBody, respHeaders] = await sendRequest(currencyServiceAddr, CURRENCY_SERVICE, GET_SUPPORTED_CURRENCIES_RPC, binReq, headers, currencyServiceTimeout);
    return unmarshalResponse(respBody, respHeaders, GET_SUPPORTED_CURRENCIES_RPC);
}

async function convert(request, headers = {}) {
    const binReq = marshalRequest(request);
    const [respBody, respHeaders] = await sendRequest(currencyServiceAddr, CURRENCY_SERVICE, CONVERT_RPC, binReq, headers, currencyServiceTimeout);
    return unmarshalResponse(respBody, respHeaders, CONVERT_RPC);
}

const currencyServiceAddr = process.env.CURRENCY_SERVICE_ADDR;
if (!currencyServiceAddr) {
    console.error("CURRENCY_SERVICE_ADDR environment variable not set");
    throw new Error("CURRENCY_SERVICE_ADDR environment variable not set");
}

let t = process.env.CURRENCY_SERVICE_TIMEOUT;
if (!t) {
    t = DEFAULT_TIMEOUT;
} else {
    t = parseInt(t, 10);
    if (isNaN(t) || t <= 0) {
        t = DEFAULT_TIMEOUT;
    }
}
const currencyServiceTimeout = t

module.exports = {getSupportedCurrencies, convert};
