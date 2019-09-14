var express = require('express');
require('dotenv').config();
var app = express();

const clientConnection = require("../gateway/fabricClient.js");
const CHANNEL = process.env['CANAL'];


app.get('/getAllGrades/', async function (req, res) {

    var requestData = {
        channel: CHANNEL,
        contractName: "grades",
        transaction: 'queryAllGrades',
        userName: 'user1',
        args: [''],
      };
    try {
    const result = await clientConnection.queryFabric(requestData);
    res.status(200).json({response: JSON.parse(result.toString())});
    //res.status(200).json({response: result});
    
    }catch (error) {
        console.error('Error al evaluar la transacción', error);
        res.status(500).json({response: 'Error al evaluar la transacción', error });
        return;
    }
});

app.post('/addGrade/', async function (req, res) {

    console.log(req.body);
    var requestData = {
        channel: CHANNEL,
        contractName: "grades",
        transaction: 'createGrade',
        userName: 'user1',
        args: [req.body.gradeId,
               req.body.studentName,
               req.body.course, 
               req.body.practice, 
               req.body.theory 
               ],
      };
      
    try {
    const result = await clientConnection.queryFabric(requestData);
    res.status(201).json({response: "Las notas han sido guardadas."});

    }catch (error) {

        res.status(500).json({response: 'Error al evaluar la transacción', error });
        return;
    }
});

app.put('/updateGrades/', async function (req, res) {

    console.log(req.body);
    var requestData = {
        channel: CHANNEL,
        contractName: 'grades',
        transaction: 'changeGrades',
        userName: 'user1',
        args: 
            [  req.body.gradeId,
               req.body.theory, 
               req.body.practice, 
            ],
      };
      
    try {
    const result = await clientConnection.queryFabric(requestData);
    res.status(200).json({response: "La nota ha sido actualizada correctamente."});

    }catch (error) {

        res.status(500).json({response: 'Error al evaluar la transacción', error });
        return;
    }
});


app.get('/getGrade/:gradeId', async function (req, res) {

    var requestData = {
        channel: CHANNEL,
        contractName: 'grades',
        transaction: 'queryGrade',
        userName: 'user1',
        args: [
            req.params.gradeId
        ],
      };
    try {
    const result = await clientConnection.queryFabric(requestData);
    res.status(200).json({response: JSON.parse(result.toString())});
    }catch (error) {
        res.status(404).json({response: 'Usuario no encontrado.'});
        return;
    }
});

app.get('/history/:gradeId', async function(req, res){
    var requestData = {
        channel: CHANNEL,
        contractName: 'grades',
        transaction: 'getHistory',
        userName: 'user1',
        args: [
            req.params.gradeId
        ]
    }
    try {
        const result = await clientConnection.queryFabric(requestData);
	console.log(result);
        res.status(200).json({ response: JSON.parse(result.toString())});
    }catch (error){
        res.status(404).json({response: 'Nota de auditoria no encontrada.'});
        return;
    }
});

module.exports = app;
