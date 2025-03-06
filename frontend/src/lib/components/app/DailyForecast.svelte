<script lang="ts">
  import type { DailyForecast } from "$type/Forecast";
  import { formatTemperature, formatPrecipitation, weatherIcon } from "$util/weather";
  import { formatDateWithWeekday } from "$util/datetime";

  const { data }: { 
    data: DailyForecast;
  } = $props<{
    data: DailyForecast;
  }>();
</script>

<div>
  {#if data !== undefined && data.time.length > 0}
    <table class="weather-table">
      <tbody>
        {#each data.time as _, index}
          <tr>
            <td>{formatDateWithWeekday(data.time[index])}</td>
            <td>
              <div class="weather-icon">
                {weatherIcon(data.data.weather_code[index])}
              </div>
            </td>
            <td class="primary">{formatTemperature(data.data.temperature_2m_max[index])}</td>
            <td class="primary">{formatTemperature(data.data.temperature_2m_min[index])}</td>
            <td class="secondary">{formatPrecipitation(data.data.precipitation_sum[index])}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Loading weather data...</p>
  {/if}
</div>
