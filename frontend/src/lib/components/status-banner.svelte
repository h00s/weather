<script lang="ts">
  import CloudOff from "@lucide/svelte/icons/cloud-off";
  import WifiOff from "@lucide/svelte/icons/wifi-off";

  type Props = { offline: boolean; stale: boolean };
  let { offline, stale }: Props = $props();
</script>

{#if offline || stale}
  <p role="status" class="panel mt-3 flex items-center gap-2 px-4 py-2 text-sm">
    {#if offline}
      <WifiOff class="size-4 shrink-0" aria-hidden="true" />Niste povezani s internetom. Prikazani su zadnji spremljeni podaci.
    {:else}
      <CloudOff class="size-4 shrink-0" aria-hidden="true" />Prognozu trenutno ne možemo osvježiti. Prikazani su zadnji podaci.
    {/if}
  </p>
{/if}
