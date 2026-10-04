<script lang="ts">
  import { formatDegrees, formatHour, formatPercent } from "$lib/helpers/format";
  import { weatherInfo } from "$lib/helpers/weather";
  import type { DailyForecast, HourlyForecast } from "$lib/types/forecast";
  import HourlyChart, { type SunEvent } from "./hourly-chart.svelte";
  import ScrollStrip from "./scroll-strip.svelte";

  type Props = { hours: HourlyForecast[]; days: DailyForecast[]; timeZone: string };
  let { hours, days, timeZone }: Props = $props();

  const sunEvents = $derived(
    days
      .flatMap((d): SunEvent[] => [
        { time: new Date(d.sunrise), kind: "sunrise" },
        { time: new Date(d.sunset), kind: "sunset" },
      ])
      .filter((e) => e.time.getUTCFullYear() > 1), // no sunrise (polar day): the zero time
  );
</script>

<section aria-labelledby="hourly-title" class="panel py-4">
  <h2 id="hourly-title" class="text-muted-foreground mb-2 px-5 text-xs font-medium tracking-wider uppercase">Idućih 24 sata</h2>
  <!-- The chart is drawn for the eye; the table below says the same to screen readers. -->
  <ScrollStrip label="Prognoza po satima" class="px-2">
    <div aria-hidden="true" class="w-max"><HourlyChart {hours} {timeZone} {sunEvents} /></div>
  </ScrollStrip>
  <!-- sr-only goes on a wrapper: a table ignores the 1 px box and the clipping, and would stretch the
       page (sideways on a phone, past the end on a desktop). -->
  <div class="sr-only">
    <table>
      <caption>Prognoza po satima</caption>
      <thead>
        <tr><th scope="col">Sat</th><th scope="col">Vrijeme</th><th scope="col">Temperatura</th><th scope="col">Vjerojatnost oborine</th></tr>
      </thead>
      <tbody>
        {#each hours as h (h.time)}
          <tr>
            <th scope="row">{formatHour(new Date(h.time), timeZone)}</th>
            <td>{weatherInfo(h.weatherCode, h.isDay).label}</td>
            <td>{formatDegrees(h.temperature)}</td>
            <td>{formatPercent(h.precipitationProbability)}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</section>
