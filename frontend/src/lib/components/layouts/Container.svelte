<script lang="ts">
  import type { Snippet } from 'svelte';

  const { children }: {
    children: Snippet;
  } = $props<{
    children: Snippet;
  }>();
</script>

<div class="flex justify-center p-4 md:p-8">
  <div class="glass-card w-full md:w-[70%] p-6 md:p-8 lg:p-10">
    {@render children()}
  </div>
</div>
