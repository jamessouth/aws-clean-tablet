import { defineConfig } from 'vite'
import reactRefresh from '@vitejs/plugin-react-refresh'
import WindiCSS from 'vite-plugin-windicss'

// import pluginTypography from "windicss/plugin/typography";
// import colors from 'windicss/colors';

export default defineConfig({
  plugins: [
    reactRefresh(),
    WindiCSS({
      preflight: false,
      scan: {
        dirs: ['./src'], // all files in the cwd
        fileExtensions: ['bs.js'], // also enabled scanning for js/ts
      },
      config: {
        // safelist: ['prose', 'prose-sm', 'm-auto'],
        darkMode: false, // or 'media' or 'class'
        // plugins: [pluginTypography],
        theme: {
          extend: {
            animation: {
              'blink': 'blink 1s infinite',
              'change': 'change 35s linear forwards 1',
              'fadein': 'fadein 4.5s ease-in forwards 1',
              'ping1': 'ping1 1s cubic-bezier(0, 0, 0.2, 1) infinite',
              'rotate': 'rotate 2.5s linear infinite',
            },
            backgroundImage: {
              'lose1': '-webkit-image-set(url("../../assets/lose11x.webp") 1x, url("../../assets/lose12x.webp") 2x)',
              'lose2': '-webkit-image-set(url("../../assets/lose21x.webp") 1x, url("../../assets/lose22x.webp") 2x)',
              'lose3': '-webkit-image-set(url("../../assets/lose31x.webp") 1x, url("../../assets/lose32x.webp") 2x)',
              'lose4': '-webkit-image-set(url("../../assets/lose41x.webp") 1x, url("../../assets/lose42x.webp") 2x)',
              'win': '-webkit-image-set(url("../../assets/win1x.webp") 1x, url("../../assets/win2x.webp") 2x)',
            },
            fontFamily: {
              'anon': 'Anonymous Pro, monospace',
              'arch': 'Architects Daughter, cursive',
              'flow': 'Indie Flower, cursive',
              'fred': 'Fredericka the Great, cursive',
              'luck': 'Luckiest Guy, cursive',
              'over': 'Overpass, sans-serif',
              'perm': 'Permanent Marker, cursive',
            },
            keyframes: {
              blink: {
                '0%, 100%': { 'transform': 'translateY(-25%)' },
                '30%': { 'opacity': '0.4' },
                '50%': { 'transform': 'translateY(0)', 'opacity': '1' },
                '79%': { 'opacity': '0.5' },
              },
              change: {
                '100%': { 'stroke-dashoffset': '1000' },
              },
              fadein: {
                'to': { 'opacity': '0.55' }
              },
              ping1: {
                '0%': {
                  'opacity': '0'
                },
                '15%, 30%': {
                  'opacity': '1'
                },
                '85%, 100%': {
                  'opacity': '0'
                },
              },
              rotate: {
                'to': { 'filter': 'hue-rotate(360deg)' },
              },
            },
            screens: {
              'newgmimg': '459px',//11/12*459=421
              'tablewidth': '550px',
              'desk': '1440px',
            },
            textShadow: {
              'lead': '0px 2px 2px #abc4d0',
              'win': '0px 0px 3px #f5f5f4',
            },
          },
        },
      },
    }),
  ],
  clearScreen: false,
  define: {
    "global": {},
  }
})
