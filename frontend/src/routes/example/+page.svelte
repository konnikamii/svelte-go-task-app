<script lang="ts">
  let { data } = $props()
</script>

<svelte:head>
  <title>Example</title>
</svelte:head>

<div>
  <h1 class="mb-4 text-2xl font-bold">Example Page</h1>
  <p>This is an example page to demonstrate layout and data passing in SvelteKit.</p>

  <div class="mt-6">
    <h2 class="text-lg font-semibold">Data from Server:</h2>
    <pre>{JSON.stringify(data, null, 2)}</pre>
  </div>
</div>
