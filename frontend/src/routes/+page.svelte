<script lang="ts">
  import Container from "$comp/layouts/Container.svelte";
  import WeatherCard from "$comp/app/WeatherCard.svelte";
  import HourlyForecast from "$comp/app/HourlyForecast.svelte";
  import DailyForecast from "$comp/app/DailyForecast.svelte";
  import type { PageData } from "./$types";

  const { data }: { data: PageData } = $props<{ data: PageData }>();
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
    <HourlyForecast data={data.forecast.hourly} />
  </div>

  <div class="pt-4">
    <DailyForecast data={data.forecast.daily} />
  </div>
</Container>
