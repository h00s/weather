<script lang="ts">
  import { onMount } from "svelte";
  import { prefersReducedMotion } from "svelte/motion";
  import AirQualityCard from "$lib/components/air-quality-card.svelte";
  import AppFooter from "$lib/components/app-footer.svelte";
  import AppHeader from "$lib/components/app-header.svelte";
  import CurrentConditions from "$lib/components/current-conditions.svelte";
  import DailyPanel from "$lib/components/daily-panel.svelte";
  import DetailTiles from "$lib/components/detail-tiles.svelte";
  import ForecastSkeleton from "$lib/components/forecast-skeleton.svelte";
  import HourlyPanel from "$lib/components/hourly-panel.svelte";
  import LocationSheet from "$lib/components/location-sheet.svelte";
  import SkyBackdrop from "$lib/components/sky-backdrop.svelte";
  import WarningBanner from "$lib/components/warning-banner.svelte";
  import Welcome from "$lib/components/welcome.svelte";
  import { todayOf, upcomingHours } from "$lib/helpers/forecast";
  import { locationLabel } from "$lib/helpers/location";
  import { skyAt } from "$lib/helpers/sky";
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
  const hint = $derived(forecast ? rainHint(forecast.hourly, now, forecast.timezone) : "");
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
  <LocationSheet bind:open={sheetOpen} />
{/if}
