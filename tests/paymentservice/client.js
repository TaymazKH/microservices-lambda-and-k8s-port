const axios = require('axios');
const {ChargeResponse} = require('./genproto/demo_pb');
const {status} = require('@grpc/grpc-js');

const DEFAULT_TIMEOUT = 10;

const PAYMENT_SERVICE = "payment-service";
const CHARGE_RPC = "charge";

function determineMessageType(rpcName) {
    return ChargeResponse;
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
        throw {code: grpcStatus, message: String(respBody)};
    }
}

module.exports = {
    marshalRequest,
    unmarshalResponse,
    sendRequest,
    DEFAULT_TIMEOUT,
    PAYMENT_SERVICE,
    CHARGE_RPC
};
