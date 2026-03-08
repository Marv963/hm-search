// @ts-check
import { defineConfig } from "astro/config";
import tailwindcss from "@tailwindcss/vite";

// https://astro.build/config
export default defineConfig({
  site: "https://marv963.github.io",
  base: "/hm-search",
  vite: {
    plugins: [tailwindcss()],
  },
});

