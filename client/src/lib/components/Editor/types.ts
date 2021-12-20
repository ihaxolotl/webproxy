export const enum KeyCode {
    Backspace = 8,
    Tab = 9,
    Enter = 13,
    Escape = 27,
    Parenteses = 57,
    M = 77,
    Y = 89,
    Z = 90,
    BackQuote = 192,
    Brackets = 219,
    Quote = 222,
};

export type Record = {
    value: string,
    selectionStart: number,
    selectionEnd: number,
};

export type History = {
    stack: Array<Record & { timestamp: number }>,
    offset: number,
};

