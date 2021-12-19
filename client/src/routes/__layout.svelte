<script lang="ts">
    import Header from "$lib/Header.svelte";

    let title: string = "My Cool Project";
</script>

<Header title={title} />
<div class="layout">
    <slot />
</div>

<style>
    .layout {
        margin: 0 136px;
        padding: 24px 0px;
    }

</style>
