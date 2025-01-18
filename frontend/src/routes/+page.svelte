<script>
  import City from '$comp/weather/City.svelte';
  import Grid from '$comp/layout/Grid.svelte';
  // import Card from '$comp/layout/Card.svelte';
</script>

<Grid cols={1}>
  <City cityName="New York" temperature={24} condition="Partly Cloudy" timestamp="Tuesday, 10:30 AM" />
</Grid>
