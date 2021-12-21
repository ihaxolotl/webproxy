<script lang="ts">
    import * as Svelte from "svelte";
    import { getLines } from "./util";
    import type { Record, History } from "./types";
    import { KeyCode } from "./types";
    import "./theme.css";

    function DeadFunction(_: any) {}

    // CSS id selector of the textarea element
    export let id: string = "codeArea";

    // Input name of the textarea
    export let name: string = "editor";

    // Form owner of the editor
    export let form: string = "";

    // Default value of the textarea
    export let value: string = "";

    // Placeholder value of the textarea
    export let placeholder: string = "";

    // Flag for autofocus
    export let isAutoFocus: boolean = false;

    // Flag for disabling the editor
    export let isDisabled: boolean = false 

    // Flag for setting the editor's data as a required field in a containing form.
    export let isRequired: boolean = false;

    // Flag for setting the editor's data as read-only.
    export let isReadOnly: boolean = false;

    // Event handlers
    export let highlight = DeadFunction;
    export let onBlur = DeadFunction;
    export let onClick = DeadFunction; 
    export let onFocus = DeadFunction;
    export let onKeyDown = DeadFunction;
    export let onKeyUp = DeadFunction;
    export let onValueChange = DeadFunction;

    // Default config variables
    const tabSize: number = 2;
    const insertSpaces: boolean = true;
    const ignoreTabKey: boolean = false;
    const HISTORY_LIMIT = 100;
    const HISTORY_TIME_GAP = 3000;
    const capture: boolean = true;

    // Reference to the textarea element
    let _input = null;

    // Editor history
    let _history: History = {
        stack: [],
        offset: -1,
    };

    // saveCurrentState saves the current state of the editor.
    function saveCurrentState() {
        if (!_input) {
            return;
        }

        updateState({
            value: _input.value,
            selectionStart: _input.selectionStart,
            selectionEnd: _input.selectionEnd,
        }, false);
    }

    function updateState(record: Record, overwrite: boolean) {
        const { stack, offset } = _history;

        if (stack.length && offset > -1) {
            // When something updates, drop the redo operations
            _history.stack = stack.slice(0, offset + 1);

            // Limit the number of operations to 100
            const count = _history.stack.length;

            if (count > HISTORY_LIMIT) {
                const extras = count - HISTORY_LIMIT;

                _history.stack = stack.slice(extras, count);
                _history.offset = Math.max(_history.offset - extras, 0);
            }
        }

        const timestamp: number = Date.now();

        if (overwrite) {
            const last = _history.stack[_history.offset];
            const reLastWord: RegExp = /[^a-z0-9]([a-z0-9]+)$/i;

            // Check if a previous entry exists and was added in a short time interval
            if (last && timestamp - last.timestamp < HISTORY_TIME_GAP) {
                const previous = getLines(last.value, last.selectionStart)
                    .pop()
                    .match(reLastWord);

                const current = getLines(record.value, record.selectionStart)
                    .pop()
                    .match(reLastWord);
                

                // If the last word of the previous line and current line matches,
                // overwrite previous entry so that undo will remove whole word
                if (previous && current && current[1].startsWith(previous[1])) {
                    _history.stack[_history.offset] = { ...record, timestamp };
                    return;
                }
            }
        }

        _history.stack.push({ ...record, timestamp });
        _history.offset++;
    }

    // updateInput sets the current state of the editor by updating the
    // values of the input element reference.
    function updateInput(record: Record) {
        if (!_input) {
            return;
        }

        _input.value = record.value;
        _input.selectionStart = record.selectionStart;
        _input.selectionEnd = record.selectionEnd;

        if (onValueChange) {
            onValueChange(record.value);
        }
    };

    // applyEdits applies the changes to the last selection state
    function _applyEdits(record: Record) {
        const last = _history.stack[_history.offset];

        if (last && _input) {
            _history.stack[_history.offset] = {
                ...last,
                selectionStart: _input.selectionStart,
                selectionEnd: _input.selectionEnd,
            };
        }

        updateState(record, false);
        updateInput(record);
    }

    // undoEdit set the editor state to the previous saved history state.
    function undoEdit() {
        const { stack, offset } = _history;
        const record: Record = stack[offset - 1];

        // If a previous edit is found, apply the changes and update the offset
        if (record) {
            updateInput(record);
            _history.offset = Math.max(offset - 1, 0);
        }
    }

    // redoEdit sets the edit at the index after the current offset of the stack
    // as the current state of the editor.
    function redoEdit() {
        const { stack, offset } = _history;
        const record: Record = stack[offset + 1];

        if (record) {
            updateInput(record);
            _history.offset = Math.min(offset + 1, stack.length - 1);
        }
    }

    // deleteSelection deletes the characters in a selection from input state.
    function deleteSelection(e: any, record: Record) {
        const hasSelection: boolean = record.selectionStart !== record.selectionEnd;
        const textBeforeCaret: string = value.substring(0, record.selectionStart);
        const tab: string = (insertSpaces ? ' ' : '\t').repeat(tabSize);

        if (textBeforeCaret.endsWith(tab) && !hasSelection) {
            e.preventDefault();

            const updatedSelection: number = record.selectionStart - tab.length;
            const updatedValue: string  = record.value.substring(0, record.selectionStart - tab.length) + record.value.substring(record.selectionEnd);

            // Remove tab character at the cursor and update the cursor position.
            _applyEdits({
                value: updatedValue,
                selectionStart: updatedSelection,
                selectionEnd: updatedSelection,
            });
        }
    }

    // addNewLine adds a new line character into the input state.
    function addNewLine(e: any, record: Record) {
        const { value, selectionStart, selectionEnd } = record;

        // Ignore selections
        if (selectionStart === selectionEnd) {
            const line = getLines(value, selectionStart).pop();
            const matches = line.match(/^\s+/);

            if (matches && matches[0]) {
                e.preventDefault();

                // Preserve indentation on inserting a new line
                const indent = '\n' + matches[0];
                const updatedSelection = selectionStart + indent.length;

                // Insert an indentation character at the current cursor position
                _applyEdits({
                    value:
                        value.substring(0, selectionStart) + indent + value.substring(selectionEnd),
                    selectionStart: updatedSelection,
                    selectionEnd: updatedSelection,
                });
            }
        }
    }

    // handleOnKeyDown handles KeyDown events in the textarea element.
    function handleOnKeyDown(e: any) {
        if (onKeyDown) {
            onKeyDown(e);

            if (e.defaultPrevented) {
                return;
            }
        }

        if (e.keyCode === KeyCode.Escape) {
            e.target.blur();
        }

        const { value, selectionStart, selectionEnd } = e.target;
        const record: Record = { value, selectionStart, selectionEnd };
        const ctrlKeyDown: boolean = e.ctrlKey && e.altKey;

        // Handle undo and redo.
        if (ctrlKeyDown) {
            if (e.shiftKey) {
                return;
            }

            switch (e.KeyCode) {
                case KeyCode.Z:
                    e.preventDefault();
                    console.log("Undo");
                    undoEdit();
                    break;
                case KeyCode.Y:
                    e.preventDefault();
                    console.log("Redo");
                    redoEdit();
                    break;
                default:
                    return
            }
        }

        switch (e.keyCode) {
            case KeyCode.Backspace:
                console.log("Deleting a character.");
                deleteSelection(e, record);
                break;
            case KeyCode.Enter:
                console.log("Adding a new line.");
                addNewLine(e, record);
                break;
            case KeyCode.Tab:
                if (ignoreTabKey) {
                    return;
                }

                if (capture) {
                    e.preventDefault();
                }
            case KeyCode.Parenteses:
            case KeyCode.Brackets:
            case KeyCode.Quote:
            case KeyCode.BackQuote:
                // handle completion of special character sequences.
                break
            default:
                break;
        }
    }

    // handleOnInput handles Input events from the textarea element.
    function handleOnInput(e: any) {
        e.preventDefault();

        const record: Record = {
            value: e.target.value,
            selectionStart: e.target.selectionStart,
            selectionEnd: e.target.selectionEnd,
        };

        updateState(record, true);
        onValueChange(e.target.value);
    }

    $: highlighted = highlight(value);

    Svelte.onMount(() => {
        saveCurrentState();
    });

</script>

<div class="container">
    <textarea
        bind:this={_input}
        bind:value={value}
        id={id}
        class="editor-overrides"
        name={name}
        on:input={handleOnInput}
        on:keydown={handleOnKeyDown}
        on:click={onClick}
        on:keyup={onKeyUp}
        on:focus={onFocus}
        on:blur={onBlur}
        disabled={isDisabled}
        form={form}
        placeholder={placeholder}
        readonly={isReadOnly}
        required={isRequired}
        autofocus={isAutoFocus}
        autocapitalize="off"
        autocomplete="off"
        autocorrect="off"
        spellcheck={false}
        data-gramm={false}
    />

    <pre class="editor-overrides editor-highlight" aria-hidden="true">
        {@html highlighted}
    </pre>
</div>

<style>
    .container, textarea, pre {
        min-height: 600px;
        width: 100%;
        font-family: 'Roboto Mono', monospace;
        font-size: 14px;
    }

    .container {
        position: relative;
        text-align: left;
        box-sizing: border-box;
        padding: 0;
        overflow: hidden;
    }

    textarea {
        position: absolute;
        top: 0px;
        left: 0px;
        height: 100%;
        width: 100%;
        resize: none;
        /* color: inherit; */
        overflow: hidden;
        -moz-osx-font-smoothing: grayscale;
        -webkit-font-smoothing: antialiased;
        -webkit-text-fill-color: transparent;
        color: transparent !important;
        caret-color: #16161D;
    }
    
    pre {
        padding-left: 60px !important;
    }

    .editor-highlight {
        position: relative;
        pointer-events: none;
    }

    .editor-overrides {
        margin: 0;
        border: 0;
        background: none;
        box-sizing: inherit;
        display: inherit;
        font-family: inherit;
        font-size: inherit;
        font-style: inherit;
        font-variant-ligatures: inherit;
        font-weight: inherit;
        letter-spacing: inherit;
        line-height: inherit;
        tab-size: inherit;
        text-indent: inherit;
        text-rendering: inherit;
        text-transform: inherit;
        white-space: pre-wrap;
        word-break: keep-all;
        overflow-wrap: break-word;
        counter-reset: line;
    }

    .editor-overrides, textarea {
        width: 100%;
        height: 100%;
        resize: none;
    }

    #codeArea {
        outline: none;
        padding-left: 60px !important;
    }
    
    /**
     * Reset the text fill color so that placeholder is visible
     */
    .editor-overrides:empty {
        -webkit-text-fill-color: inherit !important;
    }

    :global(.editorLineNumber) {
        position: absolute;
        left: 0px;
        color: #8A8A8D;
        text-align: right;
        width: 48px;
        font-weight: 400;
        padding: 0 8px 0 0;
        border-right: 1px solid #DCDCDD;
    }

    /**
     * Hack to apply some CSS on IE10 and IE11
     */
    @media all and (-ms-high-contrast: none), (-ms-high-contrast: active) {
        .editor-overrides {
            color: transparent !important;
        }

        .editor-overrides::selection {
            background-color: #accef7 !important;
            color: transparent !important;
        }
    }
</style>
