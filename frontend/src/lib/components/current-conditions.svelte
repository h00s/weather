<script lang="ts">
  import Umbrella from "@lucide/svelte/icons/umbrella";
  import { formatDegrees } from "$lib/helpers/format";
  import type { WeatherInfo } from "$lib/helpers/weather";
  import type { DailyForecast, Forecast } from "$lib/types/forecast";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { forecast: Forecast; today: DailyForecast | undefined; info: WeatherInfo; hint: string };
  let { forecast, today, info, hint }: Props = $props();

  const wet = $derived(hint !== "" && !hint.startsWith("Bez"));
</script>

<section aria-label="Trenutno vrijeme" class="text-lift flex flex-col items-center py-2 text-center sm:flex-row sm:gap-6 sm:text-left">
  <WeatherIcon name={info.icon} size={148} class="-my-3 shrink-0" />
  <div>
    <p class="text-[5.5rem] leading-none font-extralight tabular-nums sm:text-[6.5rem]">{formatDegrees(forecast.current.temperature)}</p>
    <p class="mt-1 text-2xl font-light">{info.label}</p>
    <p class="text-muted-foreground mt-1">
      Osjećaj {formatDegrees(forecast.current.apparentTemperature)}
      {#if today}
        · <span class="sr-only">najviša</span><span aria-hidden="true">↑</span>{formatDegrees(today.temperatureMax)}
        <span class="sr-only">najniža</span><span aria-hidden="true">↓</span>{formatDegrees(today.temperatureMin)}
      {/if}
    </p>
    {#if hint}
      <p class="mt-2 flex items-center justify-center gap-1.5 sm:justify-start">
        {#if wet}<Umbrella class="size-5" aria-hidden="true" />{/if}{hint}
      </p>
    {/if}
  </div>
</section>
