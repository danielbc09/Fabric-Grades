var express = require('express');
var app = express();

const clientConnection = require("../gateway/fabricClient.js");
const CHANNEL = 'channelall';
app.get('/getAllAuditedGrades/', async function (req, res) {

    var requestData = {
        channel: CHANNEL,
        contractName: 'audit',
        transaction: 'queryAllAuditGrades',
        userName: 'user1',
        args: [''],
      };
    try {
    const result = await clientConnection.queryFabric(requestData);
    res.status(200).json({response: JSON.parse(result.toString())});
    
    }catch (error) {
        console.error('Error al evaluar la transacción', error);
        res.status(500).json({response: 'Error al evaluar la transacción', error });
        return;
    }
});

app.post('/addAuditGrade', async function (req, res) {

    console.log(req.body);
    var requestData = {
        channel: CHANNEL,
        contractName: "audit",
        transaction: 'createAuditGrade',
        userName: 'user1',
        args: [req.body.gradeId,
               req.body.studentName,
               req.body.course, 
               req.body.practice, 
               req.body.theory, 
               req.body.department,
               req.body.validated,]
      };
      
    try {
    const result = await clientConnection.queryFabric(requestData);
    res.status(201).json({response: "Las notas han sido guardadas."});

    }catch (error) {

        res.status(500).json({response: 'Error al evaluar la transacción', error });
        return;
    }
});

app.get('/getAuditGrade/:gradeId/', async function (req, res) {

    var requestData = {
        channel: CHANNEL,
        contractName: "audit",
        transaction: 'queryAuditGrade',
        userName: 'user1',
        args: [
            req.params.gradeId
        ],
      };
    try {
    const result = await clientConnection.queryFabric(requestData);
    res.status(200).json({response: JSON.parse(result.toString())});
    }catch (error) {
        console.error('Error al evaluar la transacción', error);
        res.status(404).json({response: 'Usuario no encontrado.'});
        return;
    }
});

app.put('/changeStatus/', async function (req, res) {

    var requestData = {
        channel: CHANNEL,
        contractName: 'audit',
        transaction: 'changeStatus',
        userName: 'user1',
        args: 
            [  req.body.gradeId,
               req.body.validated, 
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

app.get('/history/:gradeId', async function(req, res){
    var requestData = {
        channel: CHANNEL,
        contractName: 'audit',
        transaction: 'getHistory',
	    userName: 'user1',
        args: [
            req.params.gradeId
        ]
    }
    try {
        const result = await clientConnection.queryFabric(requestData);
	const stringResult = result.toString()
	if(result.toString() === 'null'){
		throw "Nota de auditoria no encontrada."
	}else 
	{
        res.status(200).json({ response: JSON.parse(stringResult)});
	}
    }catch (error){
        res.status(404).json({response: error});
        return;
    }
});

module.exports = app;
