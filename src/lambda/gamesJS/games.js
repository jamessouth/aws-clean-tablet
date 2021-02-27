"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");

exports.handler = (req, ctx, cb) => {

    req.Records.forEach(async (rec) => {
        if (rec.eventName === "INSERT") {
            const item = rec.dynamodb.NewImage;
            if (item.pk.S.startsWith("GAME")) {
                const apiid = process.env.CT_APIID;
                const stage = process.env.CT_STAGE;
                const endpoint =
                    "https://" +
                    apiid +
                    ".execute-api." +
                    rec.awsRegion +
                    ".amazonaws.com/" +
                    stage;

                const conn = new ApiGatewayManagementApi({
                    apiVersion: "2018-11-29",
                    endpoint,
                    region: rec.AWSRegion,
                });
                try {
                    await conn
                        .postToConnection({
                            ConnectionId: item.sk.S,
                            Data: JSON.stringify({ a: 19894, b: 74156 }),
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
