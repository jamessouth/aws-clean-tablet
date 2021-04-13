"use strict";
// const AWS = require("aws-sdk");
// const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
// const DynamoDB = require("aws-sdk/clients/dynamodb");



exports.handler = (req, ctx, cb) => {
    console.log('new lmda: ', req, ctx);
    
    cb(null, `myString`);
};
