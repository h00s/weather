<script lang="ts">
  import { page } from "$app/state";
</script>

<svelte:head><title>Pogreška {page.status} · Vrijeme</title></svelte:head>

<main class="flex min-h-dvh flex-col items-center justify-center gap-4 px-4 text-center">
  <p class="text-muted-foreground text-lg">Pogreška {page.status}</p>
  <h1 class="text-3xl font-semibold">{page.status === 404 ? "Ova stranica ne postoji." : "Nešto je pošlo po zlu."}</h1>
  <a href="/" class="panel px-5 py-2.5 font-medium">Na prognozu</a>
</main>
