"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
// const DynamoDB = require("aws-sdk/clients/dynamodb");

exports.handler = (req, ctx, cb) => {
    const endpoint = `https://${req.apiid}.execute-api.${req.region}.amazonaws.com/${req.stage}`;

    const apigw = new ApiGatewayManagementApi({
        apiVersion: "2018-11-29",
        region: req.region,
        endpoint,
    });

    try {
        (async () => {
            let si;
            let ct = 3;
            si = setInterval(() => {
                
                for (const i in req.game) {
                    await apigw
                        .postToConnection({
                            ConnectionId: req.game[i].ConnID,
                            Data: JSON.stringify({
                                type: "cd",
                                count: ct,
                            }),
                        })
                        .promise();
                }
                ct--;
            }, 1000);
            if (ct === 0) {
                clearInterval(si);
            }
        })();
    } catch (err) {
        console.log("post error: ", err);
        cb(Error(err));
    }

    cb(null, `myString`);
};
