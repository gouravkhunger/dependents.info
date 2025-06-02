import type { UserConfig } from "vite"
import tailwindcss from "@tailwindcss/vite"

export default {
  server: {
    port: 4000,
    host: "0.0.0.0",
  },
  build: {
    target: "esnext",
  },
  esbuild: {
    supported: {
      "top-level-await": true,
    }
  },
  plugins: [tailwindcss()],
} satisfies UserConfig
