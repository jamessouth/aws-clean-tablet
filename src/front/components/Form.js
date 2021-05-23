import React from 'react';
// import PropTypes from 'prop-types';
// import useFormState from '../hooks/useFormState';


const ce = React.createElement;
export default function Form({
  ANSWER_MAX_LENGTH,
  answered,
  inputText,
  onEnter,
  setInputText,
}) {




  const inputBox = useRef(null);

  const INPUT_MIN_LENGTH = 2;


  const [disableSubmit, setDisableSubmit] = useState(true);
  const [isValidInput, setIsValidInput] = useState(true);
  const [badChar, setBadChar] = useState(null);
  
  
  useEffect(() => {
    const test = inputText.match(/[^a-z '-]+/i);
    if (test) {
      setBadChar(test[0]);
    } else {
      setBadChar(null);
    }
    setIsValidInput(!test);
  }, [inputText]);
  
  useEffect(() => {
    if (answered) {
      inputBox.current.blur();
    }
  }, [answered]);


  useEffect(() => {
    setDisableSubmit((inputText.length < INPUT_MIN_LENGTH || inputText.length > ANSWER_MAX_LENGTH) || answered || !isValidInput);
  }, [inputText, ANSWER_MAX_LENGTH, answered, isValidInput]);

    return ce(
      "section",
      {
        className: "relative flex flex-col justify-between items-center h-40 text-xl mb-12"
      },
      !isValidInput && ce(
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
              onEnter();
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
            onEnter();
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