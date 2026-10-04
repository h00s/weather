<script lang="ts">
  import LocationPicker from "./location-picker.svelte";
</script>

<main class="mx-auto flex min-h-dvh w-full max-w-xl flex-col justify-center gap-8 px-4 py-10">
  <header class="text-lift text-center">
    <h1 class="text-5xl font-extralight tracking-tight sm:text-6xl">Vrijeme</h1>
    <p class="text-muted-foreground mt-3 text-lg">Prognoza za mjesto u kojem jeste: sada, po satima i za idućih 7 dana.</p>
  </header>
  <section aria-label="Odabir mjesta" class="panel p-5 sm:p-6">
    <LocationPicker />
  </section>
  <p class="text-muted-foreground text-center text-sm">
    Lokaciju pamti samo vaš uređaj, a prognozu tražimo za točku zaokruženu na kilometar.
  </p>
</main>
