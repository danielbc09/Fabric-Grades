require('dotenv').config();
const { FileSystemWallet, Gateway} = require('fabric-network');
const fs = require('fs');
const path = require('path');
///Variables de entorno.
const ARCHIVO_CONEXION = process.env['ARCHIVO_DE_CONEXION'];
const ccpPath = path.resolve(__dirname, '..' ,'..', 'grades-network', ARCHIVO_CONEXION);
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);


exports.queryFabric = async function (requestData) {
    const walletPath = path.join(process.cwd(), '/wallet');
    const wallet = new FileSystemWallet(walletPath);
    const gateway = new Gateway();
    await gateway.connect(ccp, { wallet, identity: requestData.userName, discovery: { enabled: false, asLocalhost:true } });
    const network = await gateway.getNetwork(requestData.channel);   
    const contract = network.getContract(requestData.contractName);
    return await contract.submitTransaction(requestData.transaction, ...requestData.args);
};
