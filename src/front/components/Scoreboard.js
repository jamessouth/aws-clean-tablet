import React, { useEffect, useState } from "react";
// import PropTypes from 'prop-types';
// import { ul } from '../styles/Scoreboard.module.css';
// import playerSort from '../utils/playerSort';
// import mapFn from '../utils/mapFn';

const ce = React.createElement;

export default function Scoreboard({
    playerName,
    players,
    // showAnswers,
    // winners,
    // word,
}) {
    // const [toggleFinalRoundAnswers, setToggleFinalRoundAnswers] = useState(false);

    // useEffect(() => {
    //   if (winners) {
    //     setTimeout(() => {
    //       setToggleFinalRoundAnswers(!toggleFinalRoundAnswers)
    //     }, 5000);
    //   }
    // }, [toggleFinalRoundAnswers, winners]);

    // const scoreList = players
    //   .sort(playerSort('score', -1))
    //   .map(mapFn('score'));

    // const rank = scoreList.findIndex(l => l.key.split('_')[0] == playerName) + 1;

    // const answerList = players
    //   .sort(playerSort('answer', 1))
    //   .map(mapFn('answer'));

    // const titleBegin = showAnswers || toggleFinalRoundAnswers ? 'Last word:' : 'Scores:';

    // const titleEnd = showAnswers || toggleFinalRoundAnswers ? word : `You're no. ${rank}!`;

    // return (
    //   <div style={{ height: `calc(82px + (28px * ${players.length}))`, width: '100%' }}>
    //     <h2 style={{ marginBottom: '1.25em' }}>{ titleBegin }&nbsp;{ titleEnd }</h2>
    //     <ul
    //       aria-label={ showAnswers || toggleFinalRoundAnswers ? 'answers' : 'scores' }
    //       className={ ul }
    //     >
    //       { showAnswers || toggleFinalRoundAnswers ? answerList : scoreList }
    //     </ul>
    //   </div>
    // );

    return ce(
        "div",
        {
            className: "w-full",
            style: {
                height: `calc(82px + (28px * ${players.length}))`,
            },
        },
        ce(
            "h2",
            {
                className: "mb-5",
            },
            "score"
        ),
        ce(
            "ul",
            {
                className:
                    "bg-yellow-300 opacity-80 border-2 border-solid border-yellow-400 p-3 w-11/12 max-w-lg my-0 mx-auto flex flex-col justify-around items-center",
            },
            players.map((p) =>
                ce(
                    "li",
                    {
                        className:
                            "w-full flex flex-row h-7 py-0 px-2 justify-between items-center text-xl",
                        style: {
                            backgroundColor: p.color.S,
                        },
                        key: p.connid.S,
                    },
                    ce("p", null, p.name.S),
                    ce("p", null, p.score.N)
                )
            )
        )
    );

    // <li style={{backgroundColor: pl.color}} className={ li } key={ pl.name + '_' + ind }>
    //     <p>{ pl.name }</p>
    //     <p>{ pl[crit] }</p>
    //   </li>
}

// Scoreboard.propTypes = {
//   playerName: PropTypes.string,
//   players: PropTypes.array,
//   showAnswers: PropTypes.bool,
//   winners: PropTypes.bool,
//   word: PropTypes.string
// }
