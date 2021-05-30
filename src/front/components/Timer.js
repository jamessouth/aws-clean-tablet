import React from 'react';
import PropTypes from 'prop-types';

export default function Timer({ timer }) {

  return (
    <p
      role="alert"
      style={{color: '#ffee58', fontSize: '5em'}}
    >
      { timer }
    </p>
  );

}

Timer.propTypes = {
  timer: PropTypes.number
}