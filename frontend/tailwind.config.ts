import type { Config } from 'tailwindcss';
import flowbitePlugin from 'flowbite/plugin'

export default {
  content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        primary: '#01161e',
        secondary: '#124559',
        tertiary: '#598392',

        background: '#01161e',
        surface: '#124559',

        blue: {
          900: '#01161e',
          700: '#124559',
          500: '#598392',
          300: '#aeC3b0',
          100: '#eff6e0'
        },
      }
    }
  },
  plugins: [flowbitePlugin]
} as Config;