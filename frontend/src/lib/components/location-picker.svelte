<script lang="ts">
  import History from "@lucide/svelte/icons/history";
  import LoaderCircle from "@lucide/svelte/icons/loader-circle";
  import LocateFixed from "@lucide/svelte/icons/locate-fixed";
  import Search from "@lucide/svelte/icons/search";
  import { isAbortError } from "$lib/api/client";
  import { cities } from "$lib/data/cities";
  import { searchPlaces } from "$lib/services/places";
  import { locationStore } from "$lib/stores/location.svelte";
  import type { SavedLocation } from "$lib/types/location";
  import type { Place } from "$lib/types/place";

  type Props = { onchosen?: () => void };
  let { onchosen }: Props = $props();
  const id = $props.id();

  let query = $state("");
  let results = $state<Place[]>([]);
  /** The query the results (or the failure) answer: shown only while it matches the input. */
  let answered = $state("");
  let failed = $state(false);
  let searching = $state(false);
  let input = $state<HTMLInputElement>();
  let list = $state<HTMLUListElement>();

  let timer: ReturnType<typeof setTimeout> | undefined;
  let controller: AbortController | undefined;

  const trimmed = $derived(query.trim());
  const current = $derived(answered !== "" && answered === trimmed);
  const status = $derived(
    !current ? "" : failed ? "Pretraživanje ne radi" : results.length === 1 ? "1 rezultat" : `${results.length} rezultata`,
  );

  /** Searches 250 ms after the last keystroke, cancelling the request before it. */
  function search() {
    clearTimeout(timer);
    controller?.abort();
    const q = query.trim();
    if (q.length < 2) {
      results = [];
      answered = "";
      searching = false;
      return;
    }
    searching = true;
    timer = setTimeout(async () => {
      const mine = (controller = new AbortController());
      try {
        results = await searchPlaces(q, mine.signal);
        failed = false;
      } catch (e) {
        if (isAbortError(e)) return;
        results = [];
        failed = true;
      }
      answered = q;
      searching = false;
    }, 250);
  }

  // A pending search dies with the picker (when the sheet closes).
  $effect(() => () => {
    clearTimeout(timer);
    controller?.abort();
  });

  function pick(location: SavedLocation) {
    locationStore.choose(location);
    query = "";
    results = [];
    answered = "";
    onchosen?.();
  }

  const fromSearch = (p: Place): SavedLocation => ({
    kind: "place",
    name: p.name,
    county: p.county || null,
    latitude: p.latitude,
    longitude: p.longitude,
  });

  async function useMyLocation() {
    if (await locationStore.locate()) onchosen?.();
  }

  /** Enter takes the best match; the down arrow moves into the results. */
  function onInputKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && current && results.length) {
      e.preventDefault();
      pick(fromSearch(results[0]));
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      list?.querySelector("button")?.focus();
    }
  }

  /** The arrows move between results; up from the first goes back to the input. */
  function onResultKeydown(e: KeyboardEvent) {
    const buttons = [...(list?.querySelectorAll("button") ?? [])];
    const i = buttons.indexOf(e.currentTarget as HTMLButtonElement);
    if (e.key === "ArrowDown") {
      e.preventDefault();
      buttons[Math.min(i + 1, buttons.length - 1)]?.focus();
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      (i <= 0 ? input : buttons[i - 1])?.focus();
    }
  }
</script>

<div class="flex flex-col gap-6">
  <div>
    <button
      type="button"
      onclick={useMyLocation}
      disabled={locationStore.locating}
      class="text-navy flex w-full items-center justify-center gap-2 rounded-2xl bg-white px-5 py-3.5 text-base font-semibold shadow-lg transition hover:bg-white/90 disabled:opacity-70"
    >
      {#if locationStore.locating}
        <LoaderCircle class="size-5 animate-spin" aria-hidden="true" /> Tražim vašu lokaciju…
      {:else}
        <LocateFixed class="size-5" aria-hidden="true" /> Koristi moju lokaciju
      {/if}
    </button>
    {#if locationStore.error}
      <p role="alert" class="text-warn-yellow mt-2 text-sm">{locationStore.error}</p>
    {/if}
  </div>

  <div>
    <label for="search-{id}" class="sr-only">Pretraži mjesto</label>
    <div class="relative">
      <Search class="text-muted-foreground pointer-events-none absolute top-1/2 left-3.5 size-5 -translate-y-1/2" aria-hidden="true" />
      <input
        bind:this={input}
        bind:value={query}
        oninput={search}
        onkeydown={onInputKeydown}
        id="search-{id}"
        type="search"
        placeholder="Pretraži mjesto…"
        autocomplete="off"
        spellcheck="false"
        enterkeyhint="search"
        aria-describedby="search-status-{id}"
        class="placeholder:text-muted-foreground w-full rounded-2xl bg-white/10 py-3 pr-11 pl-11 text-base focus:bg-white/15"
      />
      {#if searching}
        <LoaderCircle class="text-muted-foreground pointer-events-none absolute top-1/2 right-3.5 size-5 -translate-y-1/2 animate-spin" aria-hidden="true" />
      {/if}
    </div>
    <p id="search-status-{id}" class="sr-only" aria-live="polite">{status}</p>
    {#if current && results.length}
      <ul bind:this={list} aria-label="Rezultati pretrage" class="mt-2 flex flex-col gap-1">
        {#each results as place (`${place.name}|${place.latitude}|${place.longitude}`)}
          <li>
            <button
              type="button"
              onclick={() => pick(fromSearch(place))}
              onkeydown={onResultKeydown}
              class="flex w-full flex-col rounded-xl px-3 py-2 text-left hover:bg-white/10 focus:bg-white/15"
            >
              <span class="font-medium">{place.name}</span>
              {#if place.county}<span class="text-muted-foreground text-sm">{place.county}</span>{/if}
            </button>
          </li>
        {/each}
      </ul>
    {:else if current && !searching}
      <p class="text-muted-foreground mt-2 px-1 text-sm">
        {failed ? "Pretraživanje trenutno ne radi. Odaberite grad s popisa." : `Nema mjesta „${answered}”.`}
      </p>
    {/if}
  </div>

  {#if locationStore.recent.length}
    <section aria-labelledby="recent-{id}">
      <h2 id="recent-{id}" class="text-muted-foreground mb-2 text-xs font-medium tracking-wider uppercase">Nedavno</h2>
      <ul class="flex flex-wrap gap-2">
        {#each locationStore.recent as place (`${place.latitude},${place.longitude}`)}
          <li>
            <button type="button" onclick={() => pick(place)} class="flex items-center gap-1.5 rounded-full bg-white/10 px-3.5 py-1.5 text-sm hover:bg-white/20">
              <History class="size-4 opacity-70" aria-hidden="true" />{place.name}
            </button>
          </li>
        {/each}
      </ul>
    </section>
  {/if}

  <section aria-labelledby="cities-{id}">
    <h2 id="cities-{id}" class="text-muted-foreground mb-2 text-xs font-medium tracking-wider uppercase">Gradovi</h2>
    <ul class="grid grid-cols-2 gap-2 sm:grid-cols-3">
      {#each cities as city (city.name)}
        <li>
          <button type="button" onclick={() => pick(city)} class="w-full rounded-xl bg-white/8 px-3 py-2.5 text-left hover:bg-white/15">
            {city.name}
          </button>
        </li>
      {/each}
    </ul>
  </section>
</div>
