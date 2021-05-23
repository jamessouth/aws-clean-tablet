import React from 'react';
// import PropTypes from 'prop-types';


const ce = React.createElement;
export default function Word({
  onAnimationEnd,
  playerColor,
  word,
}) {

  const blankPos = word.startsWith('_') ? 'word, blank first' : 'word, blank last';
  const divstyle = "bg-smoke-100 relative w-80 h-36 flex items-center justify-center";

  return ce(
    "div",
    {
      className: divstyle,
    },
    ce(
      "svg",
      {
        preserveAspectRatio: "none",
        className: "overflow-visible absolute top-0 left-0 w-full h-full"
      },
      ce(
        "rect",
        {
          x: "0",
          y: "0",
          width: "100%",
          height: "100%",
          onAnimationEnd,
          style: { stroke: playerColor },
          className: "animate-change rect"
        }
      )
    ),
    ce(
      "p",
      {
        ariaLabel: blankPos,
        role: "alert",
        className: "text-smoke-700 text-4xl py-0 px-6 font-perm"
      },
      `${word}`
    )
  );

}

// Word.propTypes = {
//   onAnimationEnd: PropTypes.func,
//   playerColor: PropTypes.string,
//   showAnswers: PropTypes.bool,
//   showSVGTimer: PropTypes.bool,
//   word: PropTypes.string,
// }