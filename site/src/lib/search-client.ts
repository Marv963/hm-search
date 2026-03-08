import Fuse from "fuse.js";
import type { Option } from "@/types/option";

let fuse: Fuse<Option> | null = null;
let allOptions: Option[] = [];
let currentKeys: string[] = [];

export function initSearch(options: Option[]) {
  allOptions = options;
  currentKeys = ["name", "description"];
  createFuse();
  return fuse;
}

function createFuse() {
  if (allOptions.length === 0) return;

  fuse = new Fuse(allOptions, {
    keys: currentKeys,
    threshold: 0.1, // Strikter: 0.0 = exakt, 1.0 = alles
    includeScore: true,
    minMatchCharLength: 3, // Mindestens 3 Zeichen für Match
    ignoreLocation: true, // Überall im String suchen (nicht nur am Anfang)
    distance: 100,
    useExtendedSearch: false,
  });
}

export function search(
  query: string,
  searchInTitle: boolean,
  searchInDesc: boolean,
): Option[] {
  if (!query.trim()) return allOptions;

  const keys = [];
  if (searchInTitle) keys.push("name");
  if (searchInDesc) keys.push("description");

  // WICHTIG: Fuse neu erstellen wenn sich Keys ändern
  if (JSON.stringify(keys) !== JSON.stringify(currentKeys)) {
    currentKeys = keys;
    createFuse();
  }

  if (!fuse || keys.length === 0) return allOptions;

  const results = fuse.search(query);

  // Optional: Nach Score filtern (nur gute Treffer)
  return results
    .filter((r) => r.score !== undefined && r.score < 0.4)
    .map((r) => r.item);
}

export function getAllOptions(): Option[] {
  return allOptions;
}
