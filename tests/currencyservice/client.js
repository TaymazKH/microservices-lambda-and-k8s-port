const axios = require('axios');
const {GetSupportedCurrenciesResponse, Money} = require('./genproto/demo_pb');
const {status} = require('@grpc/grpc-js');

const DEFAULT_TIMEOUT = 10;

const CURRENCY_SERVICE = "currency-service";
const GET_SUPPORTED_CURRENCIES_RPC = "get-supported-currencies";
const CONVERT_RPC = "convert";

function determineMessageType(rpcName) {
    if (rpcName === GET_SUPPORTED_CURRENCIES_RPC) {
        return GetSupportedCurrenciesResponse;
    } else {
        return Money;
    }
}

async function sendRequest(addr, serviceName, rpcName, binReq, headers = {}, timeout = DEFAULT_TIMEOUT) {
    if (!headers) {
        headers = {};
    }
    headers['rpc-name'] = rpcName;
    headers['content-type'] = 'application/octet-stream';

    let response;
    try {
        response = await axios.post(`${addr}/${serviceName}`, binReq, {
            headers: headers,
            timeout: timeout * 1000,
            responseType: 'arraybuffer' // To handle binary data
        });
    } catch (error) {
        throw new Error(`Failed to send HTTP request: ${error.message}`);
    }

    if (response.status >= 400) {
        throw new Error(`Received non-OK response: ${response.status}`);
    }

    return [response.data, response.headers];
}

function marshalRequest(msg) {
    try {
        return msg.serializeBinary();
    } catch (error) {
        throw new Error(`Failed to marshal request: ${error.message}`);
    }
}

function unmarshalResponse(respBody, headers, rpcName) {
    let grpcStatus = headers['grpc-status'];
    if (!grpcStatus) {
        throw new Error("Missing grpc-status header");
    }

    try {
        grpcStatus = parseInt(grpcStatus, 10);
    } catch (error) {
        throw new Error(`Failed to parse grpc-status header: ${grpcStatus}`);
    }

    if (grpcStatus === status.OK) {
        const MessageType = determineMessageType(rpcName);

        try {
            return MessageType.deserializeBinary(new Uint8Array(respBody));
        } catch (error) {
            throw new Error(`Failed to unmarshal response: ${error.message}`);
        }
    } else {
        throw {code: grpcStatus, message: respBody};
    }
}

module.exports = {
    marshalRequest,
    unmarshalResponse,
    sendRequest,
    DEFAULT_TIMEOUT,
    CURRENCY_SERVICE,
    GET_SUPPORTED_CURRENCIES_RPC,
    CONVERT_RPC
};
