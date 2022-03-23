import { defineConfig } from 'vite'
import reactRefresh from '@vitejs/plugin-react-refresh'
import WindiCSS from 'vite-plugin-windicss'

// import pluginTypography from "windicss/plugin/typography";
import colors from 'windicss/colors';

export default defineConfig({
  plugins: [
    reactRefresh(),
    WindiCSS({
      // preflight: false,
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
              'change': 'change 35s linear forwards 1',
              'rotate': 'rotate 3s linear infinite',
            },
            fontFamily: {
              'anon': 'Anonymous Pro, monospace',
              'arch': 'Architects Daughter, cursive',
              'flow': 'Indie Flower, cursive',
              'fred': 'Fredericka the Great, cursive',
              'luck': 'Luckiest Guy, cursive',
              'perm': 'Permanent Marker, cursive',
            },
            keyframes: {
              change: {
                '100%': { 'stroke-dashoffset': '1000' },
              },
              rotate: {
                'to': { filter: 'hue-rotate(360deg)' },
              },
            },
            lineHeight: {
              '12rem': '12rem',
            },
            screens: {
              'newgmimg': '459px',//11/12*459=421
              // 'desk': '1440px',
            },
            textShadow: {
              '2xl': '1px 0px 4px #f5f5f4bb',
            },
          },
        },
        // shortcuts: {
        //   'pubbody': '',
        //   'privbody': '',
        // },
      },


    }),
  ],
  clearScreen: false,
  define: {
    "global": {},
  }
})
