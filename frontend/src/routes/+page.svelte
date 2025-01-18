<script lang="ts">
  import Weather from "$comp/weather/Weather.svelte";
  import Grid from "$comp/layout/Grid.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<Grid cols={1}>
  <Weather
    cityName="Daruvar"
    timestamp={data.forecast.current.time}
    temperature={data.forecast.current.data.temperature_2m}
    humidity={data.forecast.current.data.relative_humidity_2m}
    precipitation={data.forecast.current.data.precipitation}
  />
</Grid>
