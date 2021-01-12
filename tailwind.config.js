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
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
