import js from "@eslint/js";
import vitest from "@vitest/eslint-plugin";
import prettier from "eslint-plugin-prettier/recommended";
import { defineConfig } from "eslint/config";
import tseslint from "typescript-eslint";

export default defineConfig(
  {
    ignores: ["**/coverage", "**/dist", "**/linter", "**/node_modules"],
  },
  js.configs.recommended,
  tseslint.configs.recommended,
  prettier,
  {
    languageOptions: {
      parserOptions: {
        project: "./tsconfig.eslint.json",
        tsconfigRootDir: import.meta.dirname,
      },
    },
    rules: {
      camelcase: "warn",
      "prettier/prettier": "error",
    },
  },
  {
    files: ["tests/**"],
    plugins: { vitest },
    rules: {
      ...vitest.configs.recommended.rules,
    },
  },
);
