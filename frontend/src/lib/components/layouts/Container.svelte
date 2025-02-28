<script lang="ts">
  import type { Snippet } from 'svelte';
  
  const { children } = $props<{ children: Snippet }>();
</script>

<div class="mx-4 md:w-3/4 md:mx-auto">
  {@render children()}
</div>