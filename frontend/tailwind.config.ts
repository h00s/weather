import type { Config } from 'tailwindcss';
import flowbitePlugin from 'flowbite/plugin'

export default {
  content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        background: '#01161e',
        primary: '#eff6e0',
      }
    }
  },
  plugins: [flowbitePlugin]
} as Config;