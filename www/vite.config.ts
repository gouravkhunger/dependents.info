import type { UserConfig } from "vite"
import tailwindcss from "@tailwindcss/vite"

export default {
  server: {
    port: 4000,
    host: "0.0.0.0",
  },
  build: {
    target: "esnext",
    rollupOptions: {
      output: {
        assetFileNames: ({ names }) => {
          for (const name of names) {
            if (/\.css$/.test(name)) {
              return "assets/css/[name][extname]";
            }
          }
          return "assets/[name]-[hash][extname]";
        }
      }
    }
  },
  plugins: [tailwindcss()],
} satisfies UserConfig
