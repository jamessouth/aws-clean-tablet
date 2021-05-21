import React from 'react';
// import PropTypes from 'prop-types';
import { bg, inv, signin } from '../styles/Form.module.css';
import useFormState from '../hooks/useFormState';


const ce = React.createElement;
export default function Form({
  answered,
  dupeName,
  hasJoined,
  invalidInput,
  onEnter,
  playerName,
  playing,
  submitSignal,
}) {

  const {
    ANSWER_MAX_LENGTH,
    badChar,
    disableSubmit,
    inputBox,
    inputText,
    isValidInput,
   
    setInputText,
  } = useFormState(
    answered,
    hasJoined,
    invalidInput,
    onEnter,
    playing,
    submitSignal,
  );

    return ce(
      "section",
      {
        className: "relative flex flex-col justify-between items-center h-40 text-xl mb-12"
      },
      (invalidInput || !isValidInput) && ce(
        "p",
        {
          className: "absolute text-smoke-100 bg-smoke-800 font-bold w-11/12 max-w-xl",
          ariaLive: "assertive"
        },
        (badChar ? `${badChar}` : 'That input') + " is not allowed"
      ),
      ce(
        "label",
        {
          ariaLive: "assertive",
          htmlFor: "inputbox"
        },
        "Enter your answer:"
      ),
      ce(
        "input",
        {
          id: "inputbox",
          autoComplete: "off",
          autoFocus,
          ref: inputBox,
          value: inputText,
          spellCheck: "false",
          onKeyPress: ({ key }) => {
            if (key == 'Enter' && !disableSubmit) {
              onEnter(inputText.slice(0, ANSWER_MAX_LENGTH));
              setInputText('');
            }
          },
          onChange: ({ target: { value } }) => setInputText(value),
          type: "text",
          placeholder: `2 - ${ANSWER_MAX_LENGTH} letters`,
          ...(answered ? { 'readOnly': true } : {})
        }
      )
    );
  
}

// Form.propTypes = {
//   answered: PropTypes.bool,
//   dupeName: PropTypes.bool,
//   hasJoined: PropTypes.bool,
//   invalidInput: PropTypes.bool,
//   onEnter: PropTypes.func,
//   playerName: PropTypes.string,
//   playing: PropTypes.bool,
//   submitSignal: PropTypes.bool
// }