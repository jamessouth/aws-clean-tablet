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
              'erase': 'erase 3s cubic-bezier(.03,.74,.03,1) forwards 1',
              'change': 'change 35s linear forwards 1',
            },
            colors: {
              smoke: colors.trueGray
            },
            fontFamily: {
              'anon': 'Anonymous Pro, monospace',
              'arch': 'Architects Daughter, cursive',
              'flow': 'Indie Flower, cursive',
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
            textShadow: {
              '2xl': '1px 0px 4px #f5f5f4bb',
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




// {
//   extract: {
//     include: ['./**/*.html', 'src/**/*.{js}'],
//     exclude: [
//       'node_modules/**/*',
//       '.git/**/*',
//     ],
//   },

// }