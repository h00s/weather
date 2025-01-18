<script lang="ts">
  import Weather from "$comp/weather/Weather.svelte";
  import Grid from "$comp/layout/Grid.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;
  console.log(data);
</script>

<Grid cols={1}>
  <Weather
    cityName="Daruvar"
    timestamp="Tuesday, 10:30 AM"
    temperature={data.forecast.current.data.temperature_2m}
    humidity={data.forecast.current.data.relative_humidity_2m}
    precipitation={data.forecast.current.data.precipitation}
  />
</Grid>
