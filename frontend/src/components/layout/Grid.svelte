<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Grid {
		cols: number;
    children: Snippet;
	}

	let { cols, children }: Grid = $props();
</script>

<div class="container mx-auto p-4">
  <div class="grid grid-cols-1 md:grid-cols-{cols} gap-4">
    {@render children()}
  </div>
</div>