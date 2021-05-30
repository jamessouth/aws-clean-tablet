import { useState, useEffect, useReducer } from 'react';
import { initialState, reducer } from '../reducers/appState';

const ws = new WebSocket(process.env.WS);

export default function useAppState() {

  const [answered, setAnswered] = useState(false);
  const [connected, setConnected] = useState(false);
  const [dupeName, setDupeName] = useState(false);
  const [gameHasBegun, setGameHasBegun] = useState(false);
  const [hasJoined, setHasJoined] = useState(false);
  const [invalidInput, setInvalidInput] = useState(false);
  const [pingServer, setPingServer] = useState(true);
  const [playerColor, setPlayerColor] = useState(null);
  const [showReset, setShowReset] = useState(false);
  const [showStartButton, setShowStartButton] = useState(true);
  const [showStartTimer, setShowStartTimer] = useState(false);
  const [showSVGTimer, setShowSVGTimer] = useState(true);
  const [showWords, setShowWords] = useState(false);
  const [submitSignal, setSubmitSignal] = useState(false);
  const [timer, setTimer] = useState(null);
  const [winners, setWinners] = useState('');
  const [
    {
      h1Text,
      newWord,
      oldWord,
      playerName,
      players,
      showAnswers,
    },
    dispatch
  ] = useReducer(reducer, initialState);

  useEffect(() => {
    ws.addEventListener('open', function () {
      setConnected(true);
    }, false);
    
    ws.addEventListener('message', function (e) {
      const {
        message,
        player,
        players,
        time,
        winners,
        word
      } = JSON.parse(e.data);
    
      switch (true) {
      case !!message:
        switch (message) {
        case 'duplicate':
          setDupeName(true);
          break;
        case 'invalid':
          setInvalidInput(true);
          break;
        case 'progress':
          setGameHasBegun(true);
          setShowStartTimer(true);
          break;
        case 'reset':
          window.location.reload();
          break;
        default: // eslint-disable-next-line no-console
          console.log('no case for this message found: ', message);
        }
        break;
      case !!player: {
        const { color, name } = player;
        setPlayerColor(color);
        dispatch({ type: 'player', name });
        setHasJoined(true);
        break;
      }
      case !!players:
        dispatch({ type: 'players', players });
        setDupeName(false);
        break;
      case !!time:
        setShowStartTimer(true);
        setShowStartButton(false);
        setTimer(time - 1);
        break;
      case !!winners:
        dispatch({ type: 'winners', winners });
        setWinners(winners)
        dispatch({ type: 'word', word: '' });
        setTimeout(() => {
          setShowReset(true);
        }, 5000);
        break
      case !!word:
        setPingServer(false);
        setAnswered(false);
        setShowSVGTimer(true);
        dispatch({ type: 'word', word });
        setShowWords(true);
        setShowStartTimer(false);
        break;
      default: // eslint-disable-next-line no-console
        console.log('no case found: ', e.data);
      }
    }, false);
    
    ws.addEventListener('error', function (e) { // eslint-disable-next-line no-console
      console.log(e, Date.now());
    }, false);
    
    ws.addEventListener('close', function () {
      setConnected(false);
    }, false);
    
    return function cleanup() {
      ws.close(1000);
    };
  }, []);

  useEffect(() => {
    if (invalidInput) {
      setTimeout(() => {
        setInvalidInput(false);
      }, 3750);
    }
  }, [invalidInput]);

  function send(text) {
    if (!hasJoined) {
      ws.send(JSON.stringify({
        name: text,
      }));
    } else {
      setAnswered(true);
      setSubmitSignal(false);
      setShowSVGTimer(false);
      ws.send(JSON.stringify({
        answer: text,
      }));
    }
  }

  function message(msg) {
    ws.send(JSON.stringify({
      message: msg
    }));
  }

  return {
    answered,
    connected,
    dupeName,
    gameHasBegun,
    h1Text,
    hasJoined,
    invalidInput,
    message,
    newWord,
    oldWord,
    pingServer,
    playerColor,
    playerName,
    players,
    send,
    setShowStartButton,
    setSubmitSignal,
    showAnswers,
    showReset,
    showStartButton,
    showStartTimer,
    showSVGTimer,
    showWords,
    submitSignal,
    timer,
    winners
  };
  
}