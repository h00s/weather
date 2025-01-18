import type { Config } from 'tailwindcss';
import flowbitePlugin from 'flowbite/plugin'

export default {
  content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        primary: '#a9d6e5',

        background: '#012a4a',
        surface: '#2c7da0',

        blue: {
          900: '#2c7da0',
          700: '#468faf',
          500: '#61a5c2',
          300: '#89c2d9',
          100: '#a9d6e5',
        },
      }
    }
  },
  plugins: [flowbitePlugin]
} as Config;