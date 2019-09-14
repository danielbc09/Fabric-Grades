// Configuracion de variables de entorno

const dotenv = require('dotenv');
dotenv.config();
//dotenv.config({path: __dirname + '/.env'});

module.export = {
    connectionFile : process.env['ARCHIVO_DE_CONEXION'],
    channel : process.env['ARCHIVO_DE_CONEXION'],
    entity : process.env.ENTIDAD_CERTIFICADORA,
    organizaton : process.env.ORGANIZACION_MSP
};

console.log( "Gonococo", module.export )