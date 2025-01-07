<script>
  import { Alert, Button } from "flowbite-svelte";
  import { InfoCircleSolid } from "flowbite-svelte-icons";
</script>

<div class="p-8">
  <Alert color="green" dismissable>
    <InfoCircleSolid slot="icon" class="w-5 h-5" />
    A simple info alert with an
    <a
      href="/"
      class="font-semibold underline hover:text-green-800 dark:hover:text-green-900"
      >example link</a
    >
    . Give it a click if you like.
  </Alert>
</div>
