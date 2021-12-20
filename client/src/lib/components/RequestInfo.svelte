<script lang="ts">
    export let domain: string = "example.com";
    export let ipaddr: string = "255.255.255.255";
</script>

<p class="info">{`Request to: ${domain} [${ipaddr}]`}</p>

<style>
    .info {
        color: #505056;
        padding: 8px 0px 4px 0px;
        font-size: 16px;
    }
</style>
