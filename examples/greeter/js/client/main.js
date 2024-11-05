const {HelloRequest, ByeRequest} = require('./genproto/hello_pb');
const {sayHello, sayBye} = require('./greeter_service_stub');

async function main() {
    const args = require('minimist')(process.argv.slice(2));
    const name = args.name || "world";

    const helloReq = new HelloRequest();
    helloReq.setName(name);

    let helloResp;
    try {
        helloResp = await sayHello(helloReq);
        console.log(`Greeting: ${helloResp.getText()}`);
    } catch (e) {
        console.error(`Error calling SayHello RPC: ${e}`);
        return;
    }

    const byeReq = new ByeRequest();
    byeReq.setName(name);

    let byeResp;
    try {
        byeResp = await sayBye(byeReq);
        console.log(`Farewell: ${byeResp.getText()}`);
    } catch (e) {
        console.error(`Error calling SayBye RPC: ${e}`);
        return;
    }

    console.log(`${helloResp.getText()} - ${byeResp.getText()}`);
}

if (require.main === module) {
    main();
}
