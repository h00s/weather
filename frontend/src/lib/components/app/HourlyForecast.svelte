<script lang="ts">
  import type { HourlyForecast } from "$type/Forecast";
  import { formatTime } from "$util/datetime";
  import {
    formatTemperature,
    formatHumidity,
    formatPrecipitation,
    weatherIcon,
  } from "$util/weather";

  const {
    data,
  }: {
    data: HourlyForecast;
  } = $props<{
    data: HourlyForecast;
  }>();

  const forecast = $derived(processHourlyForecast(data));

  function processHourlyForecast(data: HourlyForecast): HourlyForecast {
    const now = new Date();
    const currentHour = now.getHours();

    const start = data.time.findIndex((timeStr) => {
      const timeDate = new Date(timeStr);
      return timeDate.getHours() >= currentHour;
    });
    const end = start + 12;

    return {
      time: data.time.slice(start, end),
      data: {
        precipitation: data.data.precipitation.slice(start, end),
        relative_humidity_2m: data.data.relative_humidity_2m.slice(start, end),
        temperature_2m: data.data.temperature_2m.slice(start, end),
        weather_code: data.data.weather_code.slice(start, end),
      },
    };
  }
</script>

<div>
  {#if forecast !== undefined && forecast.time.length > 0}
    <table class="weather-table">
      <tbody>
        {#each forecast.time as _, index}
          <tr>
            <td>{formatTime(forecast.time[index])}</td>
            <td>
              <div class="weather-icon">
                {weatherIcon(forecast.data.weather_code[index])}
              </div>
            </td>
            <td class="primary">
              {formatTemperature(forecast.data.temperature_2m[index])}
            </td>
            <td class="secondary">
              {formatHumidity(forecast.data.relative_humidity_2m[index])}
            </td>
            <td class="secondary">
              {formatPrecipitation(forecast.data.precipitation[index], true)}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Loading weather data...</p>
  {/if}
</div>
