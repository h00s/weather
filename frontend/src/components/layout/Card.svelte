<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Card {
    children: Snippet;
	}

	let { children }: Card = $props();
</script>

<div class="p-6 bg-surface rounded-lg">
  <!-- h2 class="text-xl font-bold text-white mb-2">Card Title</h2 -->
  {@render children()}
</div>