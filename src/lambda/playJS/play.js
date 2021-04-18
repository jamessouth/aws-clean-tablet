"use strict";
// const AWS = require("aws-sdk");
// const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
const SF = require("aws-sdk/clients/stepfunctions");

const smarn = process.env.SMARN;

// const endpoint = `https://${req.apiid}.execute-api.${req.region}.amazonaws.com/${req.stage}`;
exports.handler = (req, ctx, cb) => {
    console.log("njnjnjnj: ", req);

    const sf = new SF({
        apiVersion: "2016-11-23",
        region: req.region,
    });

    try {
        (async () => {
            const res = await sf
                .startSyncExecution({
                    stateMachineArn: smarn,
                    input: JSON.stringify({
                        arr: ["1", "2", "3", "4", "5", "6", "7", "8"],
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
