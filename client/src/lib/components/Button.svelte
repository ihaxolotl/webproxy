<script lang="ts">
    export let text: string = "Button";
    export let type: string = "button"
    export let color: string = "blue";
</script>

<button type={type} class="button {color}">{text}</button>

<style>
    .button {
        border-radius: 2px;
        min-width: 144px;
        height: 40px;
        padding: 8px;
        font-size: 16px;
        font-weight: 500;
        text-align: center;
        background-color: #ffffff;
    }
    
    .button:hover {
        cursor: pointer;
    }

    .red {
        border: 1px solid #DE3838;
        color: #DE3838;
    }

    .blue {
        border: 1px solid #388FDE;
        color: #388FDE;
    }

    .gray {
        border: 1px solid #C4C4C6;
        color: #8A8A8D;
    }
</style>
