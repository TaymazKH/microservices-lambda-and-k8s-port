const {Empty, CurrencyConversionRequest, Money} = require('./genproto/demo_pb');
const {getSupportedCurrencies, convert} = require('./currency_service_stub');

async function testGetSupportedCurrencies() {
    const req = new Empty();

    let resp;
    try {
        resp = await getSupportedCurrencies(req);
        const l = resp.getCurrencyCodesList().length;
        if (l === 33) {
            console.log("GetSupportedCurrencies test successful")
        }
        else {
            console.error(`GetSupportedCurrencies RPC returned a list of ${l} items instead of 33`)
        }
    } catch (e) {
        console.error(`Error calling GetSupportedCurrencies RPC: ${JSON.stringify(e)}`);
    }
}

async function testConvert() {
    const req = new CurrencyConversionRequest();
    const money = new Money();

    money.setUnits(1);
    money.setNanos(0);
    money.setCurrencyCode('EUR')
    req.setFrom(money);
    req.setToCode('USD')

    let resp;
    try {
        resp = await convert(req);
        const u = resp.getUnits(), n = resp.getNanos(), c = resp.getCurrencyCode();
        if (u === 1 && n === 130500000 && c === 'USD') {
            console.log("Convert test successful")
        }
        else {
            console.error(`Convert RPC returned an unexpected response: units=${u}, nanos=${n}, currency_code=${c}`)
        }
    } catch (e) {
        console.error(`Error calling Convert RPC: ${JSON.stringify(e)}`);
    }
}

if (require.main === module) {
    testGetSupportedCurrencies();
    testConvert();
}
