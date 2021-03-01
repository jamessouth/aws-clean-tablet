"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
const DynamoDB = require("aws-sdk/clients/dynamodb");

let res;

exports.handler = (req, ctx, cb) => {

    req.Records.forEach(async (rec) => {
        if (rec.eventName === "INSERT") {
            const item = rec.dynamodb.NewImage;
            if (item.pk.S.startsWith("GAME")) {
                const apiid = process.env.CT_APIID;
                const stage = process.env.CT_STAGE;
                const endpoint =
                    `https://${apiid}.execute-api.${rec.awsRegion}.amazonaws.com/${stage}`;

                const apigw = new ApiGatewayManagementApi({
                    apiVersion: "2018-11-29",
                    region: rec.AWSRegion,
                    endpoint,
                });

                const dyndb = new DynamoDB({
                    apiVersion: "2012-08-10",
                    region: rec.AWSRegion,
                });

                const dbParams = {
                    TableName: rec.eventSourceARN.split("/", 2)[1],
                    KeyConditionExpression: "pk = :gm", 
                    ExpressionAttributeValues: {
                     ":gm": {
                       S: "GAME"
                      }
                    },
                };

                try {
                    res = await dyndb
                            .query(dbParams)
                            .promise();
                } catch (err) {
                    console.log("db error: ", err);
                }
                res.type = "games";
                res.data = res.Items;
                console.log('data: ', res);

                try {
                    await apigw
                        .postToConnection({
                            ConnectionId: item.connid.S,
                            Data: JSON.stringify(res),
                        })
                        .promise();
                } catch (err) {
                    console.log("post error: ", err);
                }
            } else {
                console.log("other: ");
            }
        } else {
            console.log("keys", rec.dynamodb.Keys);
        }

        console.log("Stream record: ", JSON.stringify(rec, null, 2));
    });
    cb(null, `Successfully processed ${req.Records.length} records.`);
};
