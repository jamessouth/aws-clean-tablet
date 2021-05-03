"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
const DynamoDB = require("aws-sdk/clients/dynamodb");

let gamesResults;
let connsResults;
const apiid = process.env.CT_APIID;
const stage = process.env.CT_STAGE;

function objToArr(obj) {
    const arr = [];
    for (let p in obj) {
        arr.push(obj[p].M);
    }
    return arr;
}

exports.handler = (req, ctx, cb) => {
    req.Records.forEach(async (rec) => {
        const tableName = rec.eventSourceARN.split("/", 2)[1];
        const item = rec.dynamodb.NewImage;
        console.log("item: ", item);
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

        if (rec.eventName === "INSERT") {
            if (item.pk.S.startsWith("CONN")) {
                const gamesParams = {
                    TableName: tableName,
                    KeyConditionExpression: "pk = :gm",
                    FilterExpression: "#ST = :st",
                    ExpressionAttributeValues: {
                        ":gm": {
                            S: "GAME",
                        },
                        ":st": {
                            BOOL: false,
                        },
                    },
                    ExpressionAttributeNames: {
                        "#ST": "starting",
                    },
                };
                try {
                    gamesResults = await dyndb.query(gamesParams).promise();
                } catch (err) {
                    console.log("db error: ", err);
                }
                const payload = {
                    games: gamesResults.Items.map((g) => ({
                        no: g.sk.S,
                        ready: g.ready.BOOL,
                        leader: g.leader.S,
                        players: objToArr(g.players.M),
                    })),
                    type: "games",
                };

                // console.log("data: ", payload);
                try {
                    await apigw
                        .postToConnection({
                            ConnectionId: item.GSI1SK.S,
                            Data: JSON.stringify(payload),
                        })
                        .promise();
                } catch (err) {
                    console.log("post error: ", err);
                    cb(Error(err));
                }
            } else if (item.pk.S.startsWith("GAME")) {
                const payload = {
                    games: {
                        no: item.sk.S,
                        ready: item.ready.BOOL,
                        players: objToArr(item.players.M),
                    },
                    type: "games",
                };

                // console.log("data: ", payload);

                const connsParams = {
                    TableName: tableName,
                    IndexName: "GSI1",
                    KeyConditionExpression: "GSI1PK = :cn",
                    FilterExpression: "#PL = :f",
                    ExpressionAttributeValues: {
                        ":cn": {
                            S: "CONN",
                        },
                        ":f": {
                            BOOL: false,
                        },
                    },
                    ExpressionAttributeNames: {
                        "#PL": "playing",
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
                    cb(Error(err));
                }
            } else {
                console.log("stat insert: ", item);
            }
        } else if (rec.eventName === "MODIFY") {
            if (item.pk.S.startsWith("CONN")) {
                const payload = {
                    ingame: item.game.S,
                    leadertoken: item.sk.S + item.GSI1SK.S,
                    playing: item.playing.BOOL,
                    type: "user",
                };

                // console.log("data: ", payload);
                try {
                    await apigw
                        .postToConnection({
                            ConnectionId: item.GSI1SK.S,
                            Data: JSON.stringify(payload),
                        })
                        .promise();
                } catch (err) {
                    console.log("post error: ", err);
                    cb(Error(err));
                }
            } else if (item.pk.S.startsWith("GAME")) {
                
                if (item.loading.BOOL) {
                    const payload = {
                        game: {
                            no: item.sk.S,
                            players: objToArr(item.players.M),
                        },
                        type: "playing",
                    };
                    try {
                        payload.game.players.forEach(async ({ connid }) => {
                            await apigw
                                .postToConnection({
                                    ConnectionId: connid.S,
                                    Data: JSON.stringify(payload),
                                })
                                .promise();
                        });
                    } catch (err) {
                        console.log("post error: ", err);
                        cb(Error(err));
                    }
                } else {
                    const payload = {
                        games: {
                            no: item.sk.S,
                            ready: item.ready.BOOL,
                            starting: item.starting.BOOL,
                            leader: item.leader.S,
                            loading: item.loading.BOOL,
                            players: objToArr(item.players.M),
                        },
                        type: "games",
                    };
                    const connsParams = {
                        TableName: tableName,
                        IndexName: "GSI1",
                        KeyConditionExpression: "GSI1PK = :cn",
                        FilterExpression: "#PL = :f",
                        ExpressionAttributeValues: {
                            ":cn": {
                                S: "CONN",
                            },
                            ":f": {
                                BOOL: false,
                            },
                        },
                        ExpressionAttributeNames: {
                            "#PL": "playing",
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
                        cb(Error(err));
                    }
                }
            } else {
                console.log("stat modify: ", item);
            }
        } else {
            console.log("keys", rec.dynamodb.Keys);
        }

        // console.log("Stream record: ", JSON.stringify(rec, null, 2));
    });
    cb(null, `Successfully processed ${req.Records.length} records.`);
};
