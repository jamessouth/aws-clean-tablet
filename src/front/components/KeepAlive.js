import { useEffect } from 'react';
import PropTypes from 'prop-types';

export default function KeepAlive({ pingWS }) {

  const TIMEOUT = 50000; // heroku dyno times out after 55 seconds

  useEffect(() => {
    const ping = setInterval(() => {
      pingWS();
    }, TIMEOUT);

    return () => {
      clearInterval(ping);
    }

  }, [pingWS]);

  return null;
  
}

KeepAlive.propTypes = {
  pingWS: PropTypes.func
}