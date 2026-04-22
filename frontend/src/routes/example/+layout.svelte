<script lang="ts">
  let { data, children } = $props()
</script>

<div class="p-4">
  <div class="text-xl font-bold">I am layout.svelte</div>
  <div>and have access to the following data:</div>
  <pre>{JSON.stringify(data, null, 2)}</pre>
  <br />
  {@render children()}
</div>
