"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
const DynamoDB = require("aws-sdk/clients/dynamodb");

let res;

exports.handler = (req, ctx, cb) => {

    req.Records.forEach(async (rec) => {
        if (rec.eventName === "INSERT") {
            const item = rec.dynamodb.NewImage;
            if (item.pk.S === "CONN") {
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
                const payload = {
                    data: res.Items.map(({
                        sk,
                        pk,
                        name,
                        connid,
                    }) => ({
                        no: sk.S,
                        name: name.S,
                        conn: connid.S,
                        type: pk.S,
                    })),
                    type: "games",
                };

                console.log('data: ', payload);

                try {
                    await apigw
                        .postToConnection({
                            ConnectionId: item.sk.S,
                            Data: JSON.stringify(payload),
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
