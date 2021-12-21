export function getLines(text: string, position: number): string[] {
    return text.substring(0, position).split('\n');
}
