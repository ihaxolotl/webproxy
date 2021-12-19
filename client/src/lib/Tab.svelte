<script lang="ts">
    export let active: boolean = false;
    export let title: string = "";
</script>

<div class="tab" class:active="{active}" title={title}>
    {title} 
</div>

<style>
    .tab {
        color: #16161D;
        font-size: 16px;
        padding: 8px 2px 16px 2px;
        text-align: center;
        width: fit-content;
    }

    .tab::before {
        display: block;
        content: attr(title);
        font-weight: bold;
        height: 0;
        overflow: hidden;
        visibility: hidden;
    }

    .tab:hover {
        cursor: pointer;
        font-weight: 700;
    }

    .tab.active {
        font-weight: 700;
        border-bottom: 2px solid #16161D;
    }
</style>
