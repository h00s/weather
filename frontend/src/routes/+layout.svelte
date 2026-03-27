<script>
  let { children } = $props();
  import "../app.css";
</script>

<div class="weather-bg min-h-screen">
  {@render children()}
</div>
