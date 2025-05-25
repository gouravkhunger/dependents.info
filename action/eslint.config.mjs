import path from "node:path";
import { fileURLToPath } from "node:url";

import { FlatCompat } from "@eslint/eslintrc";
import js from "@eslint/js";
import typescriptEslint from "@typescript-eslint/eslint-plugin";
import tsParser from "@typescript-eslint/parser";
import jest from "eslint-plugin-jest";
import prettier from "eslint-plugin-prettier";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
  baseDirectory: __dirname,
  allConfig: js.configs.all,
  recommendedConfig: js.configs.recommended,
});

/** @type {import("eslint").Linter.Config} */
export default [
  {
    ignores: ["**/coverage", "**/dist", "**/linter", "**/node_modules"],
  },
  ...compat.extends(
    "eslint:recommended",
    "plugin:jest/recommended",
    "plugin:prettier/recommended",
    "plugin:@typescript-eslint/recommended",
  ),

  {
    plugins: {
      jest,
      prettier,
      "@typescript-eslint": typescriptEslint,
    },

    languageOptions: {
      parser: tsParser,
      ecmaVersion: 2023,
      sourceType: "module",

      parserOptions: {
        project: "./tsconfig.eslint.json",
        tsconfigRootDir: ".",
      },
    },

    rules: {
      camelcase: "warn",
      "prettier/prettier": "error",
    },
  },
];
