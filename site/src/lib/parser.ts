/**
 * Parse Nix documentation format to HTML
 */
export function parseDescription(text: string): string {
  if (!text) return "";

  return text
    .replace(
      /<https(\s*([^>]*))/gi,
      '<a href="https$1" class="text-primary hover:underline">&lt;https$1</a>',
    )
    .replace(/\[\]\(#opt-(\s*([^)]*))/gi, "<strong>$1</strong>")
    .replace(/\)/gi, "")
    .replace(/\{var\}`([^`]+)`/gi, "<strong>$1</strong>")
    .replace(
      /:::\s*\{\.note\}\n?([\s\S]*?)\n?:::/gi,
      '<div class="mt-2 p-3 bg-info/10 border border-info/20 rounded-md text-info-foreground text-sm">$1</div>',
    )
    .replace(/\n/g, "<br />");
}

export function escapeHtml(text: string): string {
  const div = document.createElement("div");
  div.textContent = text;
  return div.innerHTML;
}
