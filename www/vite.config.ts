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
  plugins: [tailwindcss()],
} satisfies UserConfig
