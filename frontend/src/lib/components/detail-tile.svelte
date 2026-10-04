<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    label: string;
    value: string;
    detail?: string;
    icon?: Snippet;
    /** A long value (two times) in a smaller size, so it fits the tile on a phone. */
    small?: boolean;
  };
  let { label, value, detail, icon, small = false }: Props = $props();
</script>

<div class="panel flex flex-col gap-1 p-4">
  <dt class="text-muted-foreground text-xs font-medium tracking-wider uppercase">{label}</dt>
  <dd class={["flex items-center gap-2 font-light tabular-nums", small ? "text-xl" : "text-2xl"]}>{@render icon?.()}{value}</dd>
  {#if detail}<dd class="text-muted-foreground text-sm">{detail}</dd>{/if}
</div>
