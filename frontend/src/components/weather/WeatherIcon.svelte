<script lang="ts">
  import {
    Sun,
    Cloud,
    Cloudy,
    CloudFog,
    CloudDrizzle,
    CloudRain,
    CloudSnow,
    CloudLightning,
    CloudHail,
  } from "lucide-svelte";

  const { code } = $props<number>();
  let size: number = 48;

  function getWeatherIcon(code: number) {
    if (code === 0) return Sun;
    if ([1, 2, 3].includes(code)) {
      if (code === 1) return Cloud;
      if (code === 2) return Cloud;
      return Cloudy;
    }
    if ([45, 48].includes(code)) return CloudFog;
    if ([51, 53, 55].includes(code)) return CloudDrizzle;
    if ([56, 57].includes(code)) return CloudDrizzle;
    if ([61, 63, 65].includes(code)) return CloudRain;
    if ([66, 67].includes(code)) return CloudRain;
    if ([71, 73, 75].includes(code)) return CloudSnow;
    if (code === 77) return CloudSnow;
    if ([80, 81, 82].includes(code)) return CloudRain;
    if ([85, 86].includes(code)) return CloudSnow;
    if (code === 95) return CloudLightning;
    if ([96, 99].includes(code)) return CloudHail;

    return Cloud;
  }

  const IconComponent = getWeatherIcon(code);
</script>

<IconComponent {size} class={"text-blue-100"} />
