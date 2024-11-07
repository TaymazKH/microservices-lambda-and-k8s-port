const {ChargeRequest, Money, CreditCardInfo} = require('./genproto/demo_pb');
const {charge} = require('./payment_service_stub');

async function testCharge() {
    const amount = new Money();
    amount.setCurrencyCode('EUR');
    amount.setUnits(10);
    amount.setNanos(1);

    const creditCard = new CreditCardInfo();
    creditCard.setCreditCardNumber('4111111111111111');
    creditCard.setCreditCardCvv(123);
    creditCard.setCreditCardExpirationYear(3000);
    creditCard.setCreditCardExpirationMonth(1);

    const req = new ChargeRequest();
    req.setAmount(amount);
    req.setCreditCard(creditCard);

    let resp;
    try {
        resp = await charge(req);
        const tid = resp.getTransactionId();
        console.log(`Charge test successful with transaction_id: ${tid}`)
    } catch (e) {
        console.log(e);
        console.log(String(e.message))
        console.error(`Error calling Charge RPC: ${JSON.stringify(e)}`);
    }
}

if (require.main === module) {
    testCharge();
}
