const initialState = {
  // h1Text: 'CLEAN TABLET',
  // newWord: '',
  // oldWord: '',
  // playerName: '',
  games: null,
  // players: [],
  // showAnswers: false,
};
  
function reducer(
  state, {
    type,
    games,
    // name,
    // players,
    // winners,
    // word,
  }) {
  switch (type) {
  // case 'player':
  //   return {
  //     ...state,
  //     playerName: name
  //   };
    case 'games': {
      console.log('reducer: ', state.games, games);
      if (!!state.games) {
        return {
          ...state,
          games: state.games.filter(g => !g.starting).map(g => g.no === games.no ? games : g)
        };
      }
      return {
        ...state,
        games
      };
    }
  // case 'players':
  //   return {
  //     ...state,
  //     oldWord: state.newWord,
  //     players,
  //     showAnswers: !!state.newWord
  //   };
  // case 'winners': {
  //   const win = winners.includes(state.playerName) ? 'YOU WON!!' : 'YOU LOST!!';
  //   return {
  //     ...state,
  //     h1Text: win
  //   };
  // }
  // case 'word':
  //   return {
  //     ...state,
  //     newWord: word,
  //     showAnswers: false
  //   };
    default:
      throw new Error('Reducer action type not recognized');
  }
}

export { initialState, reducer };