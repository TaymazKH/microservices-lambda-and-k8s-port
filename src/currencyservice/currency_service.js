const pino = require('pino');
const {status} = require("@grpc/grpc-js");

const {GetSupportedCurrenciesResponse, Money} = require('./genproto/demo_pb')

const logger = pino({
    name: 'currencyservice-logic',
    messageKey: 'message',
    formatters: {
        level(logLevelString, logLevelNum) {
            return {severity: logLevelString}
        }
    }
});

/**
 * Lists the supported currencies
 */
function getSupportedCurrencies(empty, headers) {
    logger.info("Getting supported currencies...");

    const data = require('./data/currency_conversion.json');

    const response = new GetSupportedCurrenciesResponse();
    response.setCurrencyCodesList(Object.keys(data))

    return response
}

/**
 * Converts between currencies
 */
function convert(request, headers) {
    try {
        const data = require('./data/currency_conversion.json');

        // Convert: from_currency --> EUR
        const from = request.getFrom();
        const euros = _carry({
            units: from.getUnits() / data[from.getCurrencyCode()],
            nanos: from.getNanos() / data[from.getCurrencyCode()]
        });

        euros.nanos = Math.round(euros.nanos);

        // Convert: EUR --> to_currency
        const toCode = request.getToCode();
        const result = _carry({
            units: euros.units * data[toCode],
            nanos: euros.nanos * data[toCode]
        });

        const response = new Money();
        response.setUnits(Math.floor(result.units));
        response.setNanos(Math.floor(result.nanos));
        response.setCurrencyCode(toCode);

        logger.info("conversion request successful");
        return response
    } catch (err) {
        logger.error(`conversion request failed: ${err}`);
        throw {code: status.INVALID_ARGUMENT, message: err.message};
    }
}

/**
 * Helper function that handles decimal/fractional carrying
 */
function _carry(amount) {
    const fractionSize = Math.pow(10, 9);
    amount.nanos += (amount.units % 1) * fractionSize;
    amount.units = Math.floor(amount.units) + Math.floor(amount.nanos / fractionSize);
    amount.nanos = amount.nanos % fractionSize;
    return amount;
}

module.exports = {getSupportedCurrencies, convert};
