/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

require('dotenv').config();
const FabricCAServices = require('fabric-ca-client');
const { FileSystemWallet, X509WalletMixin } = require('fabric-network');
const fs = require('fs');
const path = require('path');

///Variables de entorno.
const ARCHIVO_CONEXION = process.env['ARCHIVO_DE_CONEXION'];
const ENTIDAD_CERTIFICADORA = process.env['ENTIDAD_CERTIFICADORA'];
const ORGANIZACION_MSP = process.env['ORGANIZACION_MSP'];

const ccpPath = path.resolve(__dirname, '..' ,'..', 'fabric-samples', 'grades-network', ARCHIVO_CONEXION);
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);


async function main() {
    try {

       // console.log("Variables de conexion", ARCHIVO_CONEXION);
        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities[ENTIDAD_CERTIFICADORA].url;
        const ca = new FabricCAServices(caURL);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = new FileSystemWallet(walletPath);
        console.log(`Billetera creada en la ruta: ${walletPath}`);

        // Check to see if we've already enrolled the admin user.
        const adminExists = await wallet.exists('admin');
        if (adminExists) {
            console.log('An identity for the admin user "admin" already exists in the wallet');
            return;
        }

        // Enroll the admin user, and import the new identity into the wallet.
        const enrollment = await ca.enroll({ enrollmentID: 'admin', enrollmentSecret: 'adminpw' });
        const identity = X509WalletMixin.createIdentity(ORGANIZACION_MSP, enrollment.certificate, enrollment.key.toBytes());
        wallet.import('admin', identity);
        console.log('El ausuario Administrador ha sido registrado correctamente.');

    } catch (error) {
        console.error(`No se pudo registrar el usuario administrador debido a: ${error}`);
        process.exit(1);
    }
}

main();
