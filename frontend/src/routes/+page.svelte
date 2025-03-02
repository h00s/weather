<script lang="ts">
  import Container from "$comp/layouts/Container.svelte";
  import CurrentForecast from "$comp/app/CurrentForecast.svelte";
  import HourlyForecast from "$comp/app/HourlyForecast.svelte";
  import DailyForecast from "$comp/app/DailyForecast.svelte";
  import type { PageData } from "./$types";

  const { data }: { data: PageData } = $props<{ data: PageData }>();
</script>

<Container>
  <div class="pt-4">
    <CurrentForecast city="Daruvar" data={data.forecast.current}/>
  </div>

  <div class="pt-4">
    <HourlyForecast data={data.forecast.hourly} />
  </div>

  <div class="pt-4">
    <DailyForecast data={data.forecast.daily} />
  </div>
</Container>
