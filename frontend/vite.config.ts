import tailwindcss from "@tailwindcss/vite";
import adapter from "@sveltejs/adapter-static";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit({
      // SPA: every route falls back to the client shell, which the Go server serves.
      adapter: adapter({ fallback: "index.html" }),
      // A phone keeps the app open for days and rarely navigates: poll for a new deploy, and the
      // root layout reloads into it the next time the page is shown.
      version: { pollInterval: 15 * 60_000 },
      compilerOptions: {
        // Runes everywhere except in libraries (the default in Svelte 6).
        runes: ({ filename }) => (filename.split(/[/\\]/).includes("node_modules") ? undefined : true),
      },
    }),
  ],
  test: {
    include: ["src/**/*.test.ts"],
  },
  server: {
    // Same origin in development too: the browser talks only to Vite, which forwards /api to the
    // Go server. The object form keeps Host intact, which the backend's csrf middleware checks.
    proxy: {
      "/api": { target: "http://localhost:3000" },
    },
  },
});
