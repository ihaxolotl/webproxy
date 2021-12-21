<script lang="ts">
    import H1 from "$lib/components/H1.svelte";

    const pageTitle: string = "History";
</script>

<svelte:head>
    <title>WebProxy | {pageTitle}</title>
</svelte:head>

<H1 text={pageTitle} />
