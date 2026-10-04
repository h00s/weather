<script lang="ts">
  import X from "@lucide/svelte/icons/x";
  import LocationPicker from "./location-picker.svelte";

  type Props = { open: boolean };
  let { open = $bindable() }: Props = $props();
  let dialog = $state<HTMLDialogElement>();

  // The <dialog> owns modality: the focus trap, Escape, and the inert page behind it. open mirrors it.
  $effect(() => {
    if (!dialog) return;
    if (open && !dialog.open) dialog.showModal();
    else if (!open && dialog.open) dialog.close();
  });
</script>

<!-- A click on the backdrop lands on the dialog itself and closes it; Escape is the keyboard way. -->
<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_noninteractive_element_interactions -->
<dialog
  bind:this={dialog}
  onclose={() => (open = false)}
  onclick={(e) => {
    if (e.target === dialog) dialog?.close();
  }}
  aria-labelledby="sheet-title"
  class="bg-sheet text-foreground border-0 p-0 shadow-2xl"
>
  <div class="flex flex-col gap-5 p-5 sm:p-6">
    <div class="flex items-center justify-between gap-4">
      <h2 id="sheet-title" class="text-lg font-semibold">Odaberite mjesto</h2>
      <button type="button" aria-label="Zatvori" onclick={() => dialog?.close()} class="-mr-2 rounded-full p-2 hover:bg-white/10">
        <X class="size-5" aria-hidden="true" />
      </button>
    </div>
    {#if open}<LocationPicker onchosen={() => dialog?.close()} />{/if}
  </div>
</dialog>

<style>
  /* A bottom sheet on phones… */
  dialog {
    margin: auto 0 0;
    width: 100%;
    max-width: 100%;
    max-height: 88dvh;
    overflow-y: auto;
    border-radius: 1.5rem 1.5rem 0 0;
    padding-bottom: env(safe-area-inset-bottom);
  }
  /* …and a centred card from 640px. */
  @media (min-width: 640px) {
    dialog {
      margin: auto;
      max-width: 32rem;
      border-radius: 1.5rem;
    }
  }
  dialog::backdrop {
    background: oklch(0 0 0 / 55%);
  }
</style>
