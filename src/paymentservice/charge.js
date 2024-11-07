const cardValidator = require('simple-card-validator');
const {v4: uuidv4} = require('uuid');
const {status} = require('@grpc/grpc-js');
const pino = require('pino');

const logger = pino({
    name: 'paymentservice-charge',
    messageKey: 'message',
    formatters: {
        level(logLevelString, logLevelNum) {
            return {severity: logLevelString}
        }
    }
});

class CreditCardError extends Error {
    constructor(message) {
        super(message);
        this.code = status.INVALID_ARGUMENT;
    }
}

class InvalidCreditCard extends CreditCardError {
    constructor() {
        super(`Credit card info is invalid`);
    }
}

class UnacceptedCreditCard extends CreditCardError {
    constructor(cardType) {
        super(`Sorry, we cannot process ${cardType} credit cards. Only VISA or MasterCard is accepted.`);
    }
}

class ExpiredCreditCard extends CreditCardError {
    constructor(number, month, year) {
        super(`Your credit card (ending ${number.substr(-4)}) expired on ${month}/${year}`);
    }
}

module.exports = function charge(request) {
    const amount = request.getAmount();
    const creditCard = request.getCreditCard();
    const cardNumber = creditCard.getCreditCardNumber();
    const cardInfo = cardValidator(cardNumber);
    const {card_type: cardType, valid} = cardInfo.getCardDetails();

    if (!valid) {
        throw new InvalidCreditCard();
    }

    // Only VISA and mastercard is accepted, other card types (AMEX, dinersclub) will
    // throw UnacceptedCreditCard error.
    if (!(cardType === 'visa' || cardType === 'mastercard')) {
        throw new UnacceptedCreditCard(cardType);
    }

    // Also validate expiration is > today.
    const currentMonth = new Date().getMonth() + 1;
    const currentYear = new Date().getFullYear();
    const year = creditCard.getCreditCardExpirationYear();
    const month = creditCard.getCreditCardExpirationMonth();
    if ((currentYear * 12 + currentMonth) > (year * 12 + month)) {
        throw new ExpiredCreditCard(cardNumber.replace('-', ''), month, year);
    }

    logger.info(`Transaction processed: ${cardType} ending ${cardNumber.substr(-4)} Amount: ${amount.getCurrencyCode()}${amount.getUnits()}.${amount.getNanos()}`);

    return uuidv4();
};