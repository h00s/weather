<script>
  import { Alert, Button } from "flowbite-svelte";
  import { InfoCircleSolid } from "flowbite-svelte-icons";
  import { fly } from "svelte/transition";
</script>

<div class="p-8">
  <Alert color="green" dismissable>
    <InfoCircleSolid slot="icon" class="w-5 h-5" />
    Flowbite enabled and ready to use!  
  </Alert>
</div>
