<script lang="ts">
  import { onMount } from "svelte";
  import { prefersReducedMotion } from "svelte/motion";
  import { online } from "svelte/reactivity/window";
  import { page } from "$app/state";
  import AirQualityCard from "$lib/components/air-quality-card.svelte";
  import AppFooter from "$lib/components/app-footer.svelte";
  import AppHeader from "$lib/components/app-header.svelte";
  import CurrentConditions from "$lib/components/current-conditions.svelte";
  import DailyPanel from "$lib/components/daily-panel.svelte";
  import DetailTiles from "$lib/components/detail-tiles.svelte";
  import ForecastSkeleton from "$lib/components/forecast-skeleton.svelte";
  import HourlyPanel from "$lib/components/hourly-panel.svelte";
  import LocationSheet from "$lib/components/location-sheet.svelte";
  import PullToRefresh from "$lib/components/pull-to-refresh.svelte";
  import SkyBackdrop from "$lib/components/sky-backdrop.svelte";
  import StatusBanner from "$lib/components/status-banner.svelte";
  import WarningBanner from "$lib/components/warning-banner.svelte";
  import Welcome from "$lib/components/welcome.svelte";
  import { todayOf, upcomingHours } from "$lib/helpers/forecast";
  import { locationLabel } from "$lib/helpers/location";
  import { parsePreview, previewOffset } from "$lib/helpers/preview";
  import { oklchToHex, skyAt } from "$lib/helpers/sky";
  import { rainHint, weatherInfo } from "$lib/helpers/weather";
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

  let sheetOpen = $state(false);

  // Design preview (?at=19:40&code=61&cloud=90), read once: the page at another time or in other weather.
  const preview = parsePreview(page.url.searchParams);
  const previewing = preview.at !== undefined || preview.code !== undefined || preview.cloud !== undefined;

  // The sky fades over a minute as the day moves on, but appears at once on the first paint and in
  // a preview, instead of fading in from navy.
  let painted = $state(false);
  onMount(() => {
    // A returning GPS reader whose permission stands is followed without a prompt.
    void locationStore.relocateIfGranted();
    requestAnimationFrame(() => requestAnimationFrame(() => (painted = true)));
  });

  const forecast = $derived(weatherStore.forecast);
  const now = $derived(
    preview.at !== undefined && forecast
      ? new Date(clock.now.getTime() + previewOffset(preview.at, clock.now, forecast.timezone))
      : clock.now,
  );
  const today = $derived(forecast ? todayOf(forecast, now) : undefined);
  const sunrise = $derived(today ? new Date(today.sunrise) : undefined);
  const sunset = $derived(today ? new Date(today.sunset) : undefined);
  const isDay = $derived(sunrise && sunset ? now >= sunrise && now < sunset : true);
  const info = $derived(weatherInfo(preview.code ?? forecast?.current.weatherCode ?? 0, isDay));
  const hint = $derived(forecast ? rainHint(forecast.hourly, now, forecast.timezone) : "");
  // Before the first forecast (the welcome screen) the sky assumes an ordinary day.
  const sky = $derived.by(() => {
    const rise = sunrise ?? new Date(new Date(now).setHours(6, 30, 0, 0));
    const set = sunset ?? new Date(new Date(now).setHours(18, 30, 0, 0));
    return skyAt(now, rise, set, preview.cloud ?? forecast?.current.cloudCover ?? 0, info.precipitation);
  });

  // The browser's own bars take the top of the sky.
  $effect(() => {
    document.querySelector('meta[name="theme-color"]')?.setAttribute("content", oklchToHex(sky.top));
  });

  async function refresh() {
    if (locationStore.current?.kind === "gps") await locationStore.locate();
    await weatherStore.refresh();
  }

  const title = $derived(locationStore.current ? `${locationLabel(locationStore.current)} · Vrijeme` : "Vrijeme");
</script>

<svelte:head><title>{title}</title></svelte:head>

<SkyBackdrop
  {sky}
  precipitation={prefersReducedMotion.current ? "none" : info.precipitation}
  transitionMs={painted && !previewing ? 60_000 : 0}
/>

{#if !locationStore.current}
  <Welcome />
{:else}
  <PullToRefresh onrefresh={refresh}>
    <main class="mx-auto w-full max-w-[1100px] px-4 pt-[max(1rem,env(safe-area-inset-top))] pb-[max(2.5rem,env(safe-area-inset-bottom))] sm:px-6">
      <AppHeader
        location={locationStore.current}
        updatedAt={weatherStore.updatedAt}
        loading={weatherStore.loading}
        now={clock.now}
        onlocation={() => (sheetOpen = true)}
        onrefresh={refresh}
      />
      <StatusBanner offline={online.current === false} stale={weatherStore.stale} />

      {#if forecast && weatherStore.warnings?.length}
        <div class="mt-4"><WarningBanner warnings={weatherStore.warnings} {now} timeZone={forecast.timezone} /></div>
      {/if}

      {#if forecast}
        <div class="mt-2 grid gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,26rem)] lg:items-start">
          <div class="flex min-w-0 flex-col gap-4">
            <CurrentConditions {forecast} {today} {info} {hint} />
            <HourlyPanel hours={upcomingHours(forecast.hourly, now, 24)} days={forecast.daily} timeZone={forecast.timezone} />
            <DetailTiles {forecast} {today} />
          </div>
          <div class="flex min-w-0 flex-col gap-4">
            <DailyPanel {forecast} {now} />
            {#if weatherStore.airQuality}<AirQualityCard airQuality={weatherStore.airQuality} />{/if}
          </div>
        </div>
      {:else if weatherStore.error}
        <div role="alert" class="panel mt-6 p-6 text-center">
          <p class="text-lg">Prognoza trenutno nije dostupna.</p>
          <button type="button" onclick={refresh} class="mt-4 rounded-full bg-white/15 px-5 py-2 font-medium hover:bg-white/25">
            Pokušaj ponovno
          </button>
        </div>
      {:else}
        <ForecastSkeleton />
      {/if}

      <AppFooter />
    </main>
  </PullToRefresh>
  <LocationSheet bind:open={sheetOpen} />
{/if}
