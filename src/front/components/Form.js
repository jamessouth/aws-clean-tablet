import React from 'react';
// import PropTypes from 'prop-types';
// import useFormState from '../hooks/useFormState';


const ce = React.createElement;
export default function Form({
  answered,
  hasJoined,
  invalidInput,
  onEnter,
  playing,
  submitSignal,
}) {

  // const {
  //   ANSWER_MAX_LENGTH,
  //   badChar,
  //   disableSubmit,
  //   inputBox,
  //   inputText,
  //   isValidInput,
   
  //   setInputText,
  // } = useFormState(
  //   answered,
  //   hasJoined,
  //   invalidInput,
  //   onEnter,
  //   playing,
  //   submitSignal,
  // );

  const inputBox = useRef(null);
  const regex = /[^a-z '-]+/i;
  
  const INPUT_MIN_LENGTH = 2;
  
  const [maxLength, setMaxLength] = useState(ANSWER_MAX_LENGTH);
  const [disableSubmit, setDisableSubmit] = useState(true);
  const [isValidInput, setIsValidInput] = useState(true);
  const [badChar, setBadChar] = useState(null);

    
  useEffect(() => {
    const test = inputText.match(regex);
    if (test) {
      setBadChar(test[0]);
    } else {
      setBadChar(null);
    }
    setIsValidInput(!test);
  }, [inputText, regex]);
  
  useEffect(() => {
    if (answered) {
      inputBox.current.blur();
    }
  }, [answered]);

  
  useEffect(() => {
    if (hasJoined) {
      setMaxLength(ANSWER_MAX_LENGTH);
    }
  }, [hasJoined]);
  
  useEffect(() => {
    if (submitSignal) {
      onEnter(inputText.slice(0, ANSWER_MAX_LENGTH));
      setInputText('');
    }
  }, [inputText, onEnter, submitSignal]);

  
  
  useEffect(() => {
    setDisableSubmit((inputText.length < INPUT_MIN_LENGTH || inputText.length > maxLength) || answered || invalidInput || !isValidInput || (hasJoined && !playing));
  }, [inputText, maxLength, answered, invalidInput, isValidInput, hasJoined, playing]);

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
          className: "h-7 w-3/5 text-xl pl-1 text-left bg-transparent border-none text-smoke-700",
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
      ),
      ce(
        "button",
        {
          className: "text-2xl text-smoke-700 bg-smoke-100 h-7 w-2/3 max-w-max cursor-pointer border-none disabled:cursor-not-allowed disabled:contrast-50",
          type: "button",
          onClick: () => {
            onEnter(inputText.slice(0, ANSWER_MAX_LENGTH));
            setInputText('');
          },
          ...(disableSubmit ? { 'disabled': true } : {})
        },
        "Submit"
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