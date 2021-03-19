const colors = require('tailwindcss/colors');

module.exports = {
  purge: [],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      colors: {
        smoke: colors.trueGray
      },
      fontFamily: {
        'arch': 'Architects Daughter, cursive',
        'fred': 'Fredericka the Great, cursive',
        'luck': 'Luckiest Guy, cursive',
        'perm': 'Permanent Marker, cursive',
      },
      gridTemplateRows: {
        'gamebox': 'repeat(6, minmax(1.5rem, 1fr))',
      },
      height: {
        // '40vh': '40vh',
      },
      lineHeight: {
        '12rem': '12rem',
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
