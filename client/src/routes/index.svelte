<script lang="ts">
    import TabGroup from "$lib/TabGroup.svelte";
    import Intercept from "$lib/containers/intercept.svelte";
    import History from "$lib/containers/history.svelte";
    import Options from "$lib/containers/options.svelte";

    let tabs = [
        {title: "Intercept", active: true, view: Intercept},
        {title: "History", active: false, view: History},
        {title: "Options", active: false, view: Options},
    ];

    function setTabActive(title: string): void {
        tabs = tabs.map((t) => {
            if (t.title === title) {
                return { ...t, active: true};
            }

            return { ...t, active: false};
        });
    }
</script>

<TabGroup {tabs} onSelect={setTabActive} />

<div class="view">
    {#each tabs as t}
        {#if t.active}
            <svelte:component this={t.view} />
        {/if}
    {/each}
</div>

<style>
    .view {
        padding: 24px 0 0 0;
    }
</style>
