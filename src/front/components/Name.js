import React from 'react';
import PropTypes from 'prop-types';
import { name } from '../styles/Name.module.css';

export default function Name({ playerName }) {

  return (
    <>
      {
        playerName &&
                    <p className={ name }>
                      { playerName }
                    </p>
      }
    </>
  );

}

Name.propTypes = {
  playerName: PropTypes.string
}