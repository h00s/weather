<script lang="ts" module>
  // Meteocons by Bas Milius (MIT), animated SVG. Only the icons named here are bundled.
  import clearDay from "@bybas/weather-icons/production/fill/all/clear-day.svg?url";
  import clearNight from "@bybas/weather-icons/production/fill/all/clear-night.svg?url";
  import partlyCloudyDay from "@bybas/weather-icons/production/fill/all/partly-cloudy-day.svg?url";
  import partlyCloudyNight from "@bybas/weather-icons/production/fill/all/partly-cloudy-night.svg?url";
  import overcastDay from "@bybas/weather-icons/production/fill/all/overcast-day.svg?url";
  import overcastNight from "@bybas/weather-icons/production/fill/all/overcast-night.svg?url";
  import fogDay from "@bybas/weather-icons/production/fill/all/fog-day.svg?url";
  import fogNight from "@bybas/weather-icons/production/fill/all/fog-night.svg?url";
  import partlyCloudyDayDrizzle from "@bybas/weather-icons/production/fill/all/partly-cloudy-day-drizzle.svg?url";
  import partlyCloudyNightDrizzle from "@bybas/weather-icons/production/fill/all/partly-cloudy-night-drizzle.svg?url";
  import drizzle from "@bybas/weather-icons/production/fill/all/drizzle.svg?url";
  import sleet from "@bybas/weather-icons/production/fill/all/sleet.svg?url";
  import partlyCloudyDayRain from "@bybas/weather-icons/production/fill/all/partly-cloudy-day-rain.svg?url";
  import partlyCloudyNightRain from "@bybas/weather-icons/production/fill/all/partly-cloudy-night-rain.svg?url";
  import rain from "@bybas/weather-icons/production/fill/all/rain.svg?url";
  import partlyCloudyDaySnow from "@bybas/weather-icons/production/fill/all/partly-cloudy-day-snow.svg?url";
  import partlyCloudyNightSnow from "@bybas/weather-icons/production/fill/all/partly-cloudy-night-snow.svg?url";
  import snow from "@bybas/weather-icons/production/fill/all/snow.svg?url";
  import thunderstormsDayRain from "@bybas/weather-icons/production/fill/all/thunderstorms-day-rain.svg?url";
  import thunderstormsNightRain from "@bybas/weather-icons/production/fill/all/thunderstorms-night-rain.svg?url";
  import hail from "@bybas/weather-icons/production/fill/all/hail.svg?url";
  import notAvailable from "@bybas/weather-icons/production/fill/all/not-available.svg?url";
  import sunrise from "@bybas/weather-icons/production/fill/all/sunrise.svg?url";
  import sunset from "@bybas/weather-icons/production/fill/all/sunset.svg?url";
  import humidity from "@bybas/weather-icons/production/fill/all/humidity.svg?url";
  import umbrella from "@bybas/weather-icons/production/fill/all/umbrella.svg?url";
  import thermometer from "@bybas/weather-icons/production/fill/all/thermometer.svg?url";

  const icons: Record<string, string> = {
    "clear-day": clearDay,
    "clear-night": clearNight,
    "partly-cloudy-day": partlyCloudyDay,
    "partly-cloudy-night": partlyCloudyNight,
    "overcast-day": overcastDay,
    "overcast-night": overcastNight,
    "fog-day": fogDay,
    "fog-night": fogNight,
    "partly-cloudy-day-drizzle": partlyCloudyDayDrizzle,
    "partly-cloudy-night-drizzle": partlyCloudyNightDrizzle,
    "drizzle": drizzle,
    "sleet": sleet,
    "partly-cloudy-day-rain": partlyCloudyDayRain,
    "partly-cloudy-night-rain": partlyCloudyNightRain,
    "rain": rain,
    "partly-cloudy-day-snow": partlyCloudyDaySnow,
    "partly-cloudy-night-snow": partlyCloudyNightSnow,
    "snow": snow,
    "thunderstorms-day-rain": thunderstormsDayRain,
    "thunderstorms-night-rain": thunderstormsNightRain,
    "hail": hail,
    "not-available": notAvailable,
    "sunrise": sunrise,
    "sunset": sunset,
    "humidity": humidity,
    "umbrella": umbrella,
    "thermometer": thermometer,
  };
</script>

<script lang="ts">
  import { prefersReducedMotion } from "svelte/motion";
  import { stillSvgUrl } from "$lib/helpers/svg";

  type Props = { name: string; size?: number; class?: string };
  let { name, size = 64, class: className }: Props = $props();

  const animated = $derived(icons[name] ?? icons["not-available"]);

  // Readers who asked their system for reduced motion get a still icon: an <img> plays the
  // icons' SMIL animation whatever CSS says. Until the still is ready the animated one shows.
  let still = $state<{ of: string; url: string }>();
  $effect(() => {
    if (!prefersReducedMotion.current) return;
    const of = animated;
    let live = true;
    stillSvgUrl(of)
      .then((url) => {
        if (live) still = { of, url };
      })
      .catch(() => {});
    return () => {
      live = false;
    };
  });
  const src = $derived(prefersReducedMotion.current && still?.of === animated ? still.url : animated);
</script>

<!-- Decorative: the text beside every icon says the same thing. -->
<img {src} width={size} height={size} alt="" draggable="false" class={className} />
