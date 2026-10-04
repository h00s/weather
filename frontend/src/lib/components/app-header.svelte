<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";
  import Navigation from "@lucide/svelte/icons/navigation";
  import RotateCw from "@lucide/svelte/icons/rotate-cw";
  import { formatAgo } from "$lib/helpers/format";
  import { locationLabel } from "$lib/helpers/location";
  import type { SavedLocation } from "$lib/types/location";

  type Props = {
    location: SavedLocation;
    updatedAt: Date | null;
    loading: boolean;
    now: Date;
    onlocation: () => void;
    onrefresh: () => void;
  };
  let { location, updatedAt, loading, now, onlocation, onrefresh }: Props = $props();
</script>

<header class="flex items-start justify-between gap-3">
  <button
    type="button"
    onclick={onlocation}
    aria-label="Promijeni mjesto, sada {locationLabel(location)}"
    class="group -ml-2 min-w-0 rounded-2xl px-2 py-1 text-left hover:bg-white/10"
  >
    <span class="text-lift flex items-center gap-1.5 text-2xl font-semibold">
      {#if location.kind === "gps"}<Navigation class="size-4 shrink-0 fill-current" aria-hidden="true" />{/if}
      <span class="truncate">{locationLabel(location)}</span>
      <ChevronDown class="size-5 shrink-0 opacity-70 transition group-hover:translate-y-0.5" aria-hidden="true" />
    </span>
    {#if location.county}<span class="text-muted-foreground block truncate text-sm">{location.county}</span>{/if}
  </button>
  <div class="text-muted-foreground flex shrink-0 items-center gap-1 text-sm">
    {#if updatedAt}<span>Ažurirano {formatAgo(updatedAt, now)}</span>{/if}
    <button type="button" onclick={onrefresh} disabled={loading} aria-label="Osvježi prognozu" class="rounded-full p-2 hover:bg-white/10 disabled:opacity-60">
      <RotateCw class="size-5 {loading ? 'animate-spin' : ''}" aria-hidden="true" />
    </button>
  </div>
</header>
