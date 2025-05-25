import { pathsToModuleNameMapper } from "ts-jest";
import config from "./tsconfig.json" assert { type: "json" };

/** @type {import("ts-jest").JestConfigWithTsJest} **/
export default {
  verbose: true,
  clearMocks: true,
  preset: "ts-jest",
  reporters: ["default"],
  testEnvironment: "node",
  resolver: "ts-jest-resolver",
  testMatch: ["**/*.test.ts"],
  extensionsToTreatAsEsm: [".ts"],
  moduleFileExtensions: ["ts", "js"],
  transform: { "^.+\\.ts$": ["ts-jest"] },
  testPathIgnorePatterns: ["/dist/", "/node_modules/"],
  moduleNameMapper: pathsToModuleNameMapper(config.compilerOptions.paths, {
    prefix: "<rootDir>",
  }),
};
