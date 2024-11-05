const http = require('http');
const {status} = require('@grpc/grpc-js');

const {ChargeRequest} = require('./genproto/demo_pb');
const {charge} = require('./payment_service');

const runningInLambda = process.env.RUN_LAMBDA === "1";
const defaultPort = 8080;

const CHARGE_RPC = "charge";

function callRPC(msg, reqData) {
    return charge(msg, reqData.headers);
}

function determineMessageType(rpcName) {
    switch (rpcName) {
        case CHARGE_RPC:
            return ChargeRequest;
        default:
            return null;
    }
}

class RequestData {
    constructor({body, headers, isBase64Encoded}) {
        this.body = body;
        this.headers = headers;
        this.isBase64Encoded = isBase64Encoded;
    }
}

function decodeRequest(reqData) {
    let binReqBody;
    if (reqData.isBase64Encoded) {
        binReqBody = new Uint8Array(Buffer.from(reqData.body, 'base64'));
    } else {
        binReqBody = new Uint8Array(reqData.body);
    }

    const rpcName = reqData.headers['rpc-name'];
    const MessageType = determineMessageType(rpcName);
    if (!MessageType) {
        return [null, generateErrorResponse(status.UNIMPLEMENTED, `Unknown RPC name: ${rpcName}`)];
    }

    try {
        return [MessageType.deserializeBinary(binReqBody), null];
    } catch (error) {
        return [null, generateErrorResponse(status.INVALID_ARGUMENT, String(error))];
    }
}

function encodeResponse(msg, rpcError) {
    if (rpcError) {
        return generateErrorResponse(rpcError.code, rpcError.message);
    }

    const binRespBody = msg.serializeBinary();
    return {
        statusCode: 200,
        headers: {
            'content-type': 'application/octet-stream',
            'grpc-status': `${status.OK}`
        },
        body: (runningInLambda ? Buffer.from(binRespBody).toString('base64') : binRespBody),
        isBase64Encoded: runningInLambda
    };
}

function generateErrorResponse(code, message) {
    return {
        statusCode: 200,
        headers: {
            'content-type': 'text/plain',
            'grpc-status': `${code}`
        },
        body: message,
        isBase64Encoded: false
    };
}

async function runLambda(event, context) {
    console.log("Handler started.");
    console.log("Event data:", event);

    const reqData = new RequestData(event);
    let [reqMsg, respData] = decodeRequest(reqData);

    if (!respData) {
        try {
            const respMsg = callRPC(reqMsg, reqData);
            respData = encodeResponse(respMsg, null);
        } catch (rpcError) {
            respData = encodeResponse(null, rpcError);
        }
    }

    console.log("Response:", respData);
    console.log("Handler finished.");
    return respData;
}

function runHTTPServer() {
    const requestHandler = async (req, res) => {
        if (req.method !== 'POST') {
            res.writeHead(405, {'Content-Type': 'text/plain'});
            res.end("Method Not Allowed");
            return;
        }

        let reqBody = [];

        req.on('data', chunk => reqBody.push(chunk));
        req.on('end', async () => {
            reqBody = Buffer.concat(reqBody);

            const headers = Object.fromEntries(
                Object.entries(req.headers).map(([k, v]) => [k.toLowerCase(), v])
            );

            const reqData = new RequestData({
                body: reqBody,
                headers: headers,
                isBase64Encoded: false,
            });

            try {
                let [reqMsg, respData] = decodeRequest(reqData);

                if (!respData) {
                    try {
                        const respMsg = callRPC(reqMsg, reqData);
                        respData = encodeResponse(respMsg, null);
                    } catch (rpcError) {
                        respData = encodeResponse(null, rpcError);
                    }
                }

                res.writeHead(respData.statusCode, respData.headers);
                res.end(respData.body);

            } catch (error) {
                console.error("Error handling request:", error);
                res.writeHead(500, {'Content-Type': 'text/plain'});
                res.end("Internal Server Error");
            }
        });
    };

    const port = process.env.PORT || defaultPort;
    const server = http.createServer(requestHandler);
    server.listen(port, () => console.log(`Port: ${port}`));
}

if (require.main === module) {
    if (runningInLambda) {
        console.warn("Conflict: RUN_LAMBDA=1 and module loaded as main.");
    } else {
        runHTTPServer();
    }
}

module.exports = {runLambda};
