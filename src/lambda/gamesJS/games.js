"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
const DynamoDB = require("aws-sdk/clients/dynamodb");

let gamesResults;
let connsResults;
const apiid = process.env.CT_APIID;
const stage = process.env.CT_STAGE;

exports.handler = (req, ctx, cb) => {
    req.Records.forEach(async (rec) => {
        const tableName = rec.eventSourceARN.split("/", 2)[1];
        const item = rec.dynamodb.NewImage;
        console.log('item: ', item);
        const endpoint = `https://${apiid}.execute-api.${rec.awsRegion}.amazonaws.com/${stage}`;

        const apigw = new ApiGatewayManagementApi({
            apiVersion: "2018-11-29",
            region: rec.AWSRegion,
            endpoint,
        });

        const dyndb = new DynamoDB({
            apiVersion: "2012-08-10",
            region: rec.AWSRegion,
        });

        if (rec.eventName === "INSERT" || (rec.eventName === "MODIFY" && item.pk.S.startsWith("GAME"))) {


            const gamesParams = {
                TableName: tableName,
                KeyConditionExpression: "pk = :gm",
                ExpressionAttributeValues: {
                    ":gm": {
                        S: "GAME",
                    },
                },
            };
            try {
                gamesResults = await dyndb.query(gamesParams).promise();
            } catch (err) {
                console.log("db error: ", err);
            }
            const payload = {
                games: gamesResults.Items.map(g => ({
                    no: g.sk.S,
                    players: g.players && g.players.SS || [],
                })),
                type: "games",
            };

            console.log("data: ", payload);

            if (item.pk.S.startsWith("CONN")) {
                try {
                    await apigw
                        .postToConnection({
                            ConnectionId: item.GSI1SK.S,
                            Data: JSON.stringify(payload),
                        })
                        .promise();
                } catch (err) {
                    console.log("post error: ", err);
                }
            } else if (item.pk.S.startsWith("GAME")) {
                const connsParams = {
                    TableName: tableName,
                    IndexName: "GSI1",
                    KeyConditionExpression: "GSI1PK = :cn",
                    ExpressionAttributeValues: {
                        ":cn": {
                            S: "CONN",
                        },
                    },
                };
                try {
                    connsResults = await dyndb.query(connsParams).promise();
                } catch (err) {
                    console.log("db error: ", err);
                }

                try {
                    connsResults.Items.forEach(async ({ GSI1SK }) => {
                        await apigw
                            .postToConnection({
                                ConnectionId: GSI1SK.S,
                                Data: JSON.stringify(payload),
                            })
                            .promise();
                    });
                } catch (err) {
                    console.log("post error: ", err);
                }
            } else {
                console.log("other: ");
            }

        } else if (rec.eventName === "MODIFY" && item.pk.S.startsWith("CONN")) {
            const payload = {
                ingame: !!item.game.S,
                type: "user",
            };

            console.log("data: ", payload);
            try {
                await apigw
                    .postToConnection({
                        ConnectionId: item.GSI1SK.S,
                        Data: JSON.stringify(payload),
                    })
                    .promise();
            } catch (err) {
                console.log("post error: ", err);
            }



        } else {
            console.log("keys", rec.dynamodb.Keys);
        }

        console.log("Stream record: ", JSON.stringify(rec, null, 2));
    });
    cb(null, `Successfully processed ${req.Records.length} records.`);
};
