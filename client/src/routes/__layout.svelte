<script lang="ts">
    import Header from "$lib/Header.svelte";
    import TabGroup from "$lib/TabGroup.svelte";

    let title: string = "My Cool Project";

    let tabs = [
        {title: "Intercept", active: true},
        {title: "History", active: false},
        {title: "Options", active: false},
    ];
</script>

<Header title={title} />
<div class="layout">
    <TabGroup {tabs} />

    <div class="view">
        <slot />
    </div>
</div>

<style>
    .layout {
        margin: 0 136px;
        padding: 24px 0px;
    }

    .view {
        padding: 24px 0 0 0;
    }
</style>
