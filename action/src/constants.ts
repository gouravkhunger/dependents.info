export const API_BASE_URL =
  process.env.API_BASE_URL ?? "https://dependents.info";

export const MESSAGE = {
  initExtraction: (repo: string) =>
    `Extracting dependents of repository ${repo}.`,
  processedPage: (count: number, repo: string) =>
    `Processed page ${count} for repository ${repo}.`,
  maxPagesReached: (maxPages: number, repo: string) =>
    `Reached maximum pages limit of ${maxPages} for repository ${repo}.`,
  dependentsCount: (count: number, repo: string) =>
    `Found ${count} dependents for repository ${repo}.`,
  wroteFile: (filePath: string) => `File written successfully at ${filePath}.`,
  artifactUploadLog: (status: string, name: string) =>
    `Artifact upload ${status} for ${name}.`,
  DONE: "Action completed successfully.",
  FORK_DETECTED: "Forked repository detected. Skipping submission.",
} as const;

export const ERROR = {
  INVALID_REPO_FORMAT:
    "Repository must be in the format 'owner/repo' and contain only valid characters.",
  failedToWriteFile: (filePath: string, msg?: string) =>
    `Failed to write file at ${filePath}: ${msg ?? "unknown error"}.`,
  failedToFetch: (url: string, statusCode?: number) =>
    `Failed to fetch ${url}: ${statusCode ?? "unknown error"}.`,
  failedToSubmitData: (statusCode: number, message: string) =>
    `Failed to submit data: ${statusCode} - ${message}.`,
  contentTypeMismatch: (url: string, expected: string, actual?: string) =>
    `Expected Content-Type for ${url} to be '${expected}', but got ${actual ?? "'unknown'"}.`,
  readBufferNotSupported: (url: string) =>
    `Response from ${url} does not support readBodyBuffer.`,
  failedToParseTotalDependents: (repo: string) =>
    `Failed to parse total dependents for repository ${repo}.`,
} as const;
