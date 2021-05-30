import React from 'react';
import PropTypes from 'prop-types';

export default function Start({ gameHasBegun, onClick, players }) {

  const ppl = 3 - players.length == 1 ? 'player' : 'players';

  return (
    <div style={{width: '100%'}}>
      {
        !gameHasBegun &&
          <button
            aria-live="polite"
            type="button"
            onClick={ onClick }
            { ...(players.length < 3 ? { 'disabled': true } : {}) }
          >
            { players.length < 3 ? 'Need ' + (3 - players.length) + ' more ' + ppl : 'Start Game' }
          </button>
      }
    </div>
  );

}

Start.propTypes = {
  gameHasBegun: PropTypes.bool,
  onClick: PropTypes.func,
  players: PropTypes.array
}