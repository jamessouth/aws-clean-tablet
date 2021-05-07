const initialState = {
    // h1Text: 'CLEAN TABLET',
    // newWord: '',
    // oldWord: '',
    // playerName: '',
    games: null,
    // players: [],
    // showAnswers: false,
};

function processGames(arr, ni) {
    const list = [...arr];
    for (let i = 0; i < arr.length; i++) {
        if (arr[i].no === ni.no) {
            if (ni.starting) {
                return list.splice(i, 1);
            } else {
                list[i] = ni;
                return list;
            }
        }
    }
    return [ni, ...list];
}

function reducer(
    state,
    {
        type,
        games,
        // name,
        // players,
        // winners,
        // word,
    }
) {
    switch (type) {
        // case 'player':
        //   return {
        //     ...state,
        //     playerName: name
        //   };
        case "games": {
          console.log("state.games in reducer: ", state.games);
          console.log();
          console.log("game in reducer: ", games);
            if (!!state.games) {
                return {
                    ...state,
                    games: processGames(state.games, games),
                };
            }
            return {
                ...state,
                games,
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
            throw new Error("Reducer action type not recognized");
    }
}

export { initialState, reducer };
