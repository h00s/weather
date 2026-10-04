<div role="status" class="mt-2 grid animate-pulse gap-4 lg:grid-cols-[minmax(0,1fr)_minmax(0,26rem)]">
  <span class="sr-only">Učitavam prognozu…</span>
  <div class="flex flex-col gap-4">
    <div class="mx-auto my-4 h-36 w-64 rounded-3xl bg-white/10"></div>
    <div class="panel h-56"></div>
    <div class="grid grid-cols-2 gap-3 sm:grid-cols-3">
      {#each [1, 2, 3, 4, 5, 6] as n (n)}<div class="panel h-28"></div>{/each}
    </div>
  </div>
  <div class="panel h-[30rem]"></div>
</div>
