<script lang="ts">
    import H1 from "$lib/components/H1.svelte";
    import Button from "$lib/components/Button.svelte";
    import RequestInfo from "$lib/components/RequestInfo.svelte";
    import { Editor } from "$lib/components/Editor";
    import { highlight, languages } from "prismjs/components/prism-core";
    import "prismjs/components/prism-http";

    const pageTitle: string = "Intercept";
    let data: string = "GET / HTTP/1.1\r\nHost: localhost:9999\r\nUser-Agent: curl/7.80.0\r\nAccept: */*\r\nConnection: Keep-Alive\r\n\r\n";

    function hightlightWithLineNumbers(input: string, language: any) {
        return highlight(input, language)
            .split("\n")
            .map((line: string, i: number) => `<span class='editorLineNumber'>${i + 1}</span>${line}`)
            .join("\n");
    }

    function onValueChange(v: string): void {}
</script>

<svelte:head>
    <title>WebProxy | {pageTitle}</title>
</svelte:head>

<H1 text={pageTitle} />

<RequestInfo />

<div class="actions">
    <Button text="Forward" color="blue" />
    <Button text="Drop" color="red" />
    <Button text="Intercept On" color="gray" />
</div>

<div class="editor-container">
    <Editor
        value={data}
        onValueChange={onValueChange}
        highlight={code => hightlightWithLineNumbers(code, languages.http)}
    />
</div>

<style>
    .actions {
        display: flex;
        border-bottom-width: 1px;
        border-bottom-color: #C4C4C6;
        border-bottom-style: solid;
        gap: 12px;
        padding: 16px 0px;
    }

    .editor-container {
        max-height: 600px;
        overflow-y: auto;
        margin: 16px 0 0 0;
    }
</style>
