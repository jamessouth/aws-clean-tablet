"use strict";
// const AWS = require("aws-sdk");
// const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
const SF = require("aws-sdk/clients/stepfunctions");

const smarn = process.env.SMARN;


const colors = [
    "#dc2626",
    "#0c4a6e",
    "#16a34a",
    "#7c2d12",
    "#c026d3",
    "#365314",
    "#0891b2",
    "#581c87",
];


// const endpoint = `https://${req.apiid}.execute-api.${req.region}.amazonaws.com/${req.stage}`;



exports.handler = (req, ctx, cb) => {
   
    const sf = new SF({
        apiVersion: "2016-11-23",
        region: req.region,
    });

    const sfInput = Object.keys(req.game.Players).map((p, i) => ({
        id: p,
        color: colors[i],
        name: req.game.Players[p].Name,
        connid: req.game.Players[p].ConnID,
    }));

    console.log('sfip: ', sfInput);

    try {
        (async () => {
            const res = await sf
                .startSyncExecution({
                    stateMachineArn: smarn,
                    input: JSON.stringify({
                        gameno: req.game.Sk,
                        arr: sfInput,
                    }),
                    name: "sfex1",
                })
                .promise();

            console.log("sfresssssss: ", res);
        })();
    } catch (err) {
        cb(Error(err));
    }

    cb(null, `myString`);
};
// const apigw = new ApiGatewayManagementApi({
//     apiVersion: "2018-11-29",
//     region: req.region,
//     endpoint,
// });
// try {
//     (async () => {
//         for (const i in req.game) {
//             console.log('pc: ', );
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
