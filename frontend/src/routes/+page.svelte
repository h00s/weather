<script lang="ts">
  import { onMount } from "svelte";
  import { prefersReducedMotion } from "svelte/motion";
  import AppHeader from "$lib/components/app-header.svelte";
  import LocationSheet from "$lib/components/location-sheet.svelte";
  import SkyBackdrop from "$lib/components/sky-backdrop.svelte";
  import Welcome from "$lib/components/welcome.svelte";
  import { todayOf } from "$lib/helpers/forecast";
  import { formatDegrees } from "$lib/helpers/format";
  import { locationLabel } from "$lib/helpers/location";
  import { skyAt } from "$lib/helpers/sky";
  import { weatherInfo } from "$lib/helpers/weather";
  import { clock } from "$lib/stores/clock.svelte";
  import { locationStore } from "$lib/stores/location.svelte";
  import { weatherStore } from "$lib/stores/weather.svelte";

  $effect(() => clock.subscribe());
  $effect(() => weatherStore.subscribe());
  // Follow the chosen location; the store ignores a repeat of the same ~1 km cell.
  $effect(() => {
    const location = locationStore.current;
    if (location) weatherStore.show(location);
  });
  // A returning GPS reader whose permission stands is followed without a prompt.
  onMount(() => void locationStore.relocateIfGranted());

  let sheetOpen = $state(false);

  const now = $derived(clock.now);
  const forecast = $derived(weatherStore.forecast);
  const today = $derived(forecast ? todayOf(forecast, now) : undefined);
  const sunrise = $derived(today ? new Date(today.sunrise) : undefined);
  const sunset = $derived(today ? new Date(today.sunset) : undefined);
  const isDay = $derived(sunrise && sunset ? now >= sunrise && now < sunset : true);
  const info = $derived(weatherInfo(forecast?.current.weatherCode ?? 0, isDay));
  // Before the first forecast (the welcome screen) the sky assumes an ordinary day.
  const sky = $derived.by(() => {
    const rise = sunrise ?? new Date(new Date(now).setHours(6, 30, 0, 0));
    const set = sunset ?? new Date(new Date(now).setHours(18, 30, 0, 0));
    return skyAt(now, rise, set, forecast?.current.cloudCover ?? 0, info.precipitation);
  });

  async function refresh() {
    if (locationStore.current?.kind === "gps") await locationStore.locate();
    await weatherStore.refresh();
  }

  const title = $derived(locationStore.current ? `${locationLabel(locationStore.current)} · Vrijeme` : "Vrijeme");
</script>

<svelte:head><title>{title}</title></svelte:head>

<SkyBackdrop {sky} precipitation={prefersReducedMotion.current ? "none" : info.precipitation} />

{#if !locationStore.current}
  <Welcome />
{:else}
  <main class="mx-auto w-full max-w-[1100px] px-4 pt-[max(1rem,env(safe-area-inset-top))] pb-[max(2.5rem,env(safe-area-inset-bottom))] sm:px-6">
    <AppHeader
      location={locationStore.current}
      updatedAt={weatherStore.updatedAt}
      loading={weatherStore.loading}
      {now}
      onlocation={() => (sheetOpen = true)}
      onrefresh={refresh}
    />
    {#if forecast}
      <p class="text-lift mt-8 text-8xl font-extralight">{formatDegrees(forecast.current.temperature)}</p>
    {/if}
  </main>
  <LocationSheet bind:open={sheetOpen} />
{/if}
