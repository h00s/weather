/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />
import { build, files, version } from "$service-worker";

// Precaches the app shell so the installed app opens offline; the last forecast comes from the
// device's storage (stores/weather.svelte.ts). The API is never cached here.
const sw = self as unknown as ServiceWorkerGlobalScope;
const CACHE = `vrijeme-${version}`;
const ASSETS = new Set([...build, ...files]);
const SHELL = "/";

sw.addEventListener("install", (event) => {
  event.waitUntil(
    (async () => {
      const cache = await caches.open(CACHE);
      await cache.addAll([...ASSETS]);
      // spa/v2 answers the shell only to navigations, which it recognizes by Accept: text/html.
      await cache.add(new Request(SHELL, { headers: { Accept: "text/html" } }));
      await sw.skipWaiting();
    })(),
  );
});

sw.addEventListener("activate", (event) => {
  event.waitUntil(
    (async () => {
      for (const key of await caches.keys()) if (key !== CACHE) await caches.delete(key);
      await sw.clients.claim();
    })(),
  );
});

sw.addEventListener("fetch", (event) => {
  const request = event.request;
  if (request.method !== "GET") return;
  const url = new URL(request.url);
  if (url.origin !== sw.location.origin || url.pathname.startsWith("/api/")) return;

  if (ASSETS.has(url.pathname)) {
    event.respondWith(caches.match(url.pathname).then((cached) => cached ?? fetch(request)));
  } else if (request.mode === "navigate") {
    // Network first, so a deploy shows at once; the cached shell when offline.
    event.respondWith(fetch(request).catch(async () => (await caches.match(SHELL)) ?? Response.error()));
  }
});
