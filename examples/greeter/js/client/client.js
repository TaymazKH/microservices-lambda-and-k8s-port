const axios = require('axios');
const {HelloResponse, ByeResponse} = require('./genproto/hello_pb');
const {status} = require('@grpc/grpc-js');

const DEFAULT_TIMEOUT = 10;

const GREETER_SERVICE = "greeter";
const SAY_HELLO_RPC = "say-hello";
const SAY_BYE_RPC = "say-bye";

function determineMessageType(rpcName) {
    if (rpcName === SAY_HELLO_RPC) {
        return HelloResponse;
    } else {
        return ByeResponse;
    }
}

function sendRequest(addr, serviceName, rpcName, binReq, headers = {}, timeout = DEFAULT_TIMEOUT) {
    if (!headers) {
        headers = {};
    }
    headers['rpc-name'] = rpcName;
    headers['content-type'] = 'application/octet-stream';

    let response;
    try {
        response = axios.post(`${addr}/${serviceName}`, binReq, {
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
