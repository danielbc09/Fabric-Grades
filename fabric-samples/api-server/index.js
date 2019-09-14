'use strict';
require('dotenv').config();
var express = require('express');
var bodyParser = require('body-parser');
const SERVIDOR = process.env['SERVIDOR'];
// Variables express
var app = express();
app.use(bodyParser.json());

// Permitir conexion desde fuera
app.use(function(req, res, next) {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Credentials", "true");
    res.header("Access-Control-Allow-Methods", "POST, GET, DELETE, UPDATE, PUT");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, token");
    next();
});


// Variables
app.set('port', 8081);

console.log("*************************Bienvenido al servidor "+ SERVIDOR +"*************************");


var gradesRoute = require('./routes/grades.js');
var auditedGrades = require('./routes/audit.js');

app.use('/api/grades/', gradesRoute);
app.use('/api/audited/', auditedGrades);

app.listen(app.get('port'), function() {
    console.log('Express server ' + SERVIDOR + ' puerto 8081: \x1b[32m%s\x1b[0m', 'online');
});



