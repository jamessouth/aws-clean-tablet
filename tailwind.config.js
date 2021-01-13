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
        'luck': 'Luckiest Guy, cursive',
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
