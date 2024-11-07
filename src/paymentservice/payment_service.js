const {ChargeResponse} = require('./genproto/demo_pb');
const pino = require('pino');
const charge = require('./charge');

const logger = pino({
    name: 'paymentservice-logic',
    messageKey: 'message',
    formatters: {
        level(logLevelString, logLevelNum) {
            return {severity: logLevelString}
        }
    }
});

function handleCharge(request, headers) {
    try {
        logger.info(`PaymentService#Charge invoked with request ${request}`);
        const transaction_id = charge(request);

        const response = new ChargeResponse();
        response.setTransactionId(transaction_id);
        return response;

    } catch (err) {
        logger.warn(err);
        throw err;
    }
}

module.exports = {handleCharge};
