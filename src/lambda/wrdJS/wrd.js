"use strict";
// const AWS = require("aws-sdk");
const ApiGatewayManagementApi = require("aws-sdk/clients/apigatewaymanagementapi");
const SF = require("aws-sdk/clients/stepfunctions");

// const smarn = process.env.SMARN;

exports.handler = (req, ctx, cb) => {
    console.log("wrddddd", req);

    // const apigw = new ApiGatewayManagementApi({
    //     apiVersion: "2018-11-29",
    //     region: req.region,
    //     endpoint: req.endpoint,
    // });

    const sf = new SF({
        apiVersion: "2016-11-23",
        region: req.region,
    });

    const hiScore = 7;
    

    try {
        (async () => {
            // for (const i of req.conns) {
            //     console.log('ppcc: ', i);
            //     const res = await apigw
            //         .postToConnection({
            //             ConnectionId: i,
            //             Data: JSON.stringify({
            //                 type: "word",
            //                 word: req.word,
            //             }),
            //         })
            //         .promise();

            //     console.log("xcvxcvxres: ", res);

                setTimeout(() => {
                    sf.sendTaskSuccess(
                        {
                            output: JSON.stringify(hiScore),
                            taskToken: req.token,
                        },
                        (err, data) => {
                            if (err) console.log(err, err.stack);
                            else console.log("SFdataaaa", data);
                        }
                    );
                }, 4000);
            // }
        })();
    } catch (err) {
        console.log("post error: ", err);
        cb(Error(err));
    }

    cb(null, `myOtherString`);
};
