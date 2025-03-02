<script lang="ts">
  import Container from "$comp/layouts/Container.svelte";
  import WeatherCard from "$comp/app/WeatherCard.svelte";
  import WeatherHourly from "$comp/app/WeatherHourly.svelte";
  import type { PageData } from "./$types";

  const { data } = $props<{ data: PageData }>();
</script>

<Container>
  <div class="pt-4">
    <WeatherCard
      weather={{
        city: "Daruvar",
        timestamp: data.forecast.current.time,
        temperature: data.forecast.current.data.temperature_2m,
        humidity: data.forecast.current.data.relative_humidity_2m,
        precipitation: data.forecast.current.data.precipitation,
        code: data.forecast.current.data.weather_code,
      }}
    />
  </div>

  <div class="pt-4">
    <WeatherHourly hourlyForecast={data.forecast.hourly} />
  </div>
</Container>
