const colors = require('windicss/colors');
const typography = require('windicss/plugin/typography');

module.exports = {
  extract: {
    include: ['./**/*.html'],
  },
  safelist: ['prose', 'prose-sm', 'm-auto'],
  darkMode: false, // or 'media' or 'class'
  plugins: [typography],
  theme: {
    extend: {
      animation: {
        'erase': 'erase 3s cubic-bezier(.03,.74,.03,1) forwards 1',
        'change': 'change 35s linear forwards 1',
      },
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
      keyframes: {
        change: {
          '100%': { 'stroke-dashoffset': '1000' },
        },
        erase: {
          '100%': { opacity: '0' },
        },
      },
      lineHeight: {
        '12rem': '12rem',
      },
    },
  },
}
