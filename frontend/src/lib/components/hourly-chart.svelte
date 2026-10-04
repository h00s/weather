<script lang="ts" module>
  export interface SunEvent {
    time: Date;
    kind: "sunrise" | "sunset";
  }
</script>

<script lang="ts">
  import { formatClock, formatDegrees, formatHour } from "$lib/helpers/format";
  import { areaPath, extent, linePath, scale } from "$lib/helpers/sparkline";
  import { weatherInfo } from "$lib/helpers/weather";
  import type { HourlyForecast } from "$lib/types/forecast";
  import WeatherIcon from "./weather-icon.svelte";

  type Props = { hours: HourlyForecast[]; timeZone: string; sunEvents?: SunEvent[] };
  let { hours, timeZone, sunEvents = [] }: Props = $props();
  const id = $props.id();

  // One column per hour; the strip scrolls sideways inside its panel. Two charts share the time
  // axis, never two y-scales: temperature above, the chance of rain below.
  const SLOT = 56;
  const TEMP_H = 92;
  const PAD = 24; // room for the labels above the curve
  const RAIN_H = 34;
  const BAR = 20;

  const width = $derived(hours.length * SLOT);
  const cx = (i: number) => (i + 0.5) * SLOT;
  const temps = $derived(hours.map((h) => h.temperature));
  const domain = $derived.by((): [number, number] => {
    const [lo, hi] = extent(temps) ?? [0, 1];
    const grow = Math.max(0, 4 - (hi - lo)) / 2; // a flat day stays flat-looking
    return [lo - grow, hi + grow];
  });
  const y = $derived(scale(domain, [TEMP_H - PAD, PAD]));
  const wet = $derived(hours.some((h) => h.precipitationProbability >= 5));
  const rainY = scale([0, 100], [RAIN_H, 12]);

  // A bar with a rounded top and a square foot on the baseline.
  function bar(x: number, top: number) {
    const r = Math.min(4, (RAIN_H - top) / 2, BAR / 2);
    const l = x - BAR / 2;
    return `M${l},${RAIN_H}V${top + r}Q${l},${top} ${l + r},${top}H${l + BAR - r}Q${l + BAR},${top} ${l + BAR},${top + r}V${RAIN_H}Z`;
  }

  // Sunrise and sunset inside the strip, placed by time between the hour columns.
  const start = $derived(hours.length ? Date.parse(hours[0].time) : 0);
  const sunMarks = $derived(
    sunEvents
      .map((e) => ({ ...e, x: ((e.time.getTime() - start) / 3_600_000) * SLOT }))
      .filter((m) => m.x > 0 && m.x < width),
  );
</script>

<div class="relative pb-6" style:width="{width}px">
  <div class="flex">
    {#each hours as h, i (h.time)}
      <div class="flex flex-col items-center" style:width="{SLOT}px">
        <span class="text-muted-foreground text-sm">{i === 0 ? "Sada" : formatHour(new Date(h.time), timeZone)}</span>
        <WeatherIcon name={weatherInfo(h.weatherCode, h.isDay).icon} size={40} />
      </div>
    {/each}
  </div>

  <svg {width} height={TEMP_H} viewBox="0 0 {width} {TEMP_H}" class="block overflow-visible">
    <defs>
      <linearGradient id="temp-fill-{id}" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0" stop-color="white" stop-opacity="0.16" />
        <stop offset="1" stop-color="white" stop-opacity="0" />
      </linearGradient>
    </defs>
    <g transform="translate({SLOT / 2} 0)">
      <path d={areaPath(temps, width - SLOT, TEMP_H, domain, PAD)} fill="url(#temp-fill-{id})" />
      <path
        d={linePath(temps, width - SLOT, TEMP_H, domain, PAD)}
        fill="none"
        stroke="white"
        stroke-opacity="0.9"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </g>
    {#each temps as t, i (hours[i].time)}
      <text x={cx(i)} y={y(t) - 9} text-anchor="middle" class="fill-white text-[15px] font-medium">{formatDegrees(t)}</text>
    {/each}
  </svg>

  {#if wet}
    <svg {width} height={RAIN_H} viewBox="0 0 {width} {RAIN_H}" class="block overflow-visible">
      <line x1="0" x2={width} y1={RAIN_H - 0.5} y2={RAIN_H - 0.5} stroke="white" stroke-opacity="0.15" />
      {#each hours as h, i (h.time)}
        {#if h.precipitationProbability >= 5}
          <path d={bar(cx(i), rainY(h.precipitationProbability))} fill="var(--color-rain)" opacity="0.85" />
        {/if}
        {#if h.precipitationProbability >= 20}
          <text x={cx(i)} y={rainY(h.precipitationProbability) - 3} text-anchor="middle" class="fill-white/85 text-[12px]">
            {Math.round(h.precipitationProbability)} %
          </text>
        {/if}
      {/each}
    </svg>
  {/if}

  {#each sunMarks as m (m.time.getTime())}
    <div class="pointer-events-none absolute inset-y-0 border-l border-dashed border-amber-200/50" style:left="{m.x}px">
      <span class="absolute bottom-0 left-1 text-xs whitespace-nowrap text-amber-100">
        {m.kind === "sunrise" ? "Izlazak" : "Zalazak"} {formatClock(m.time, timeZone)}
      </span>
    </div>
  {/each}
</div>
