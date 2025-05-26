export const MESSAGE = {
  initExtraction: (repo: string) =>
    `Extracting dependents of repository ${repo}.`,
  processedPage: (count: number, repo: string) =>
    `Processed page ${count} for repository ${repo}.`,
  dependentsCount: (count: number, repo: string) =>
    `Found ${count} dependents for repository ${repo}.`,
  wroteFile: (filePath: string) => `File written successfully at ${filePath}.`,
  artifactUploadLog: (status: string, name: string) =>
    `Artifact upload ${status} for ${name}.`,
} as const;

export const ERROR = {
  INVALID_REPO_FORMAT:
    "Input 'repo' must be in the format 'owner/repo' and contain only valid characters.",
  failedToWriteFile: (filePath: string, msg?: string) =>
    `Failed to write file at ${filePath}: ${msg ?? "unknown error"}.`,
  failedToFetch: (url: string, statusCode?: number) =>
    `Failed to fetch ${url}: ${statusCode ?? "unknown error"}.`,
  contentTypeMismatch: (url: string, expected: string, actual?: string) =>
    `Expected Content-Type for ${url} to be '${expected}', but got ${actual ?? "'unknown'"}.`,
  readBufferNotSupported: (url: string) =>
    `Response from ${url} does not support readBodyBuffer.`,
  failedToParseTotalDependents: (repo: string) =>
    `Failed to parse total dependents for repository ${repo}.`,
} as const;
