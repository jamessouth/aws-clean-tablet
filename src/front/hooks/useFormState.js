import { useState, useEffect, useRef } from 'react';

export default function useFormState(
  answered,
  hasJoined,
  invalidInput,
  onEnter,
  playing,
  submitSignal,
) {



  

  return {
    ANSWER_MAX_LENGTH,
    badChar,
    disableSubmit,
    inputBox,
    inputText,
    isValidInput,
    setInputText,
  };

}