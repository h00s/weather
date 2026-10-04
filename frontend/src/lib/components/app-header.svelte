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

<header>
  <div class="flex items-center justify-between gap-3">
    <button
      type="button"
      onclick={onlocation}
      aria-label="Promijeni mjesto, sada {locationLabel(location)}"
      class="group text-lift -ml-2 flex min-w-0 items-center gap-1.5 rounded-2xl px-2 py-1 text-2xl font-semibold hover:bg-white/10"
    >
      {#if location.kind === "gps"}<Navigation class="size-4 shrink-0 fill-current" aria-hidden="true" />{/if}
      <span class="truncate">{locationLabel(location)}</span>
      <ChevronDown class="size-5 shrink-0 opacity-70 transition group-hover:translate-y-0.5" aria-hidden="true" />
    </button>
    <div class="text-muted-foreground flex shrink-0 items-center gap-1 text-sm">
      {#if updatedAt}<span>Ažurirano {formatAgo(updatedAt, now)}</span>{/if}
      <button type="button" onclick={onrefresh} disabled={loading} aria-label="Osvježi prognozu" class="rounded-full p-2 hover:bg-white/10 disabled:opacity-60">
        <RotateCw class="size-5 {loading ? 'animate-spin' : ''}" aria-hidden="true" />
      </button>
    </div>
  </div>
  {#if location.county}<p class="text-muted-foreground truncate text-sm">{location.county}</p>{/if}
</header>
