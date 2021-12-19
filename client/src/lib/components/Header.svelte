<script lang="ts">
    export let title: string = "Error";
</script>

<div class="header">
    <div class="inner">
        <div id="branding">Proxy</div>
        <div class="project">
            <h2 class="project-title">{title}</h2>
            <span class="bar">|</span>
            <div class="menu">
                <p>Proxy</p>
            </div>
        </div>
    </div>
</div>

<style>
    #branding {
        color: #16161D;
        font-size: 24px;
        font-weight: 700;
        text-align: center;
        text-transform: uppercase;
        width: 136px;
    }

    .header {
        width: 100%; 
        height: 56px;
        padding: 16px 0px;
        border-bottom-width: 1px;
        border-bottom-color: #C4C4C6;
        border-bottom-style: solid;
    }

    .inner {
        display: flex;
    }

    .project {
        display: flex;
        align-items: center;
    }

    .project-title {
        display: inline;
        font-size: 16px;
        font-weight: 700;
    }

    .bar {
        margin: 0 8px;
    }

    .menu {
        font-size: 16px;
        font-weight: 500;
    }
</style>
