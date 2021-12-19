<script lang="ts">
    import Tab from "./Tab.svelte";

    export let tabs = [];
</script>

<div class="group">
    {#each tabs as t}
        <Tab active={t.active} title={t.title} />
    {/each}
</div>

<style>
    .group {
        display: flex;
        gap: 24px;
    }
</style>
