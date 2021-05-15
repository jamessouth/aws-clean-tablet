"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
// const SF = require("aws-sdk/clients/stepfunctions");

// const smarn = process.env.SMARN;








exports.handler = (req, ctx, cb) => {
    console.log('wrddddd', req, "ccccctx", ctx);
    // const endpoint = `https://${req.apiid}.execute-api.${req.region}.amazonaws.com/${req.stage}`;


    // const apigw = new ApiGatewayManagementApi({
    //     apiVersion: "2018-11-29",
    //     region: req.region,
    //     endpoint,
    // });
    // try {
    //     (async () => {
    //         for (const i in req.game) {
    //             // console.log('pc: ', );
    //             await apigw
    //                 .postToConnection({
    //                     ConnectionId: req.game[i].ConnID,
    //                     Data: JSON.stringify({
    //                         type: "cd",
    //                         count: 5,
    //                     }),
    //                 })
    //                 .promise();
    //         }
    //     })();
    // } catch (err) {
    //     console.log("post error: ", err);
    //     cb(Error(err));
    // }


    cb(null, `myOtherString`);
};

