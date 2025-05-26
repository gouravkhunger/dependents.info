export const MESSAGE = {
  initExtraction: (name: string) =>
    `Extracting dependents of repository ${name}`,
  processedPage: (count: number, name: string) =>
    `Processed page ${count} for repository ${name}`,
  REACHED_MAX_PAGES:
    "Reached the maximum number of pages (5). Stopping further processing.",
  processedDependents: (count: number, name: string) =>
    `Processed ${count} public dependents for repository ${name}`,
} as const;

export const ERROR = {
  INVALID_REPO_FORMAT:
    "Input 'repo' must be in the format 'owner/repo' and contain only valid characters.",
  failedToFetch: (url: string, statusCode?: number) =>
    `Failed to fetch ${url}: ${statusCode ?? "unknown error"}`,
  contentTypeMismatch: (url: string, expected: string, actual?: string) =>
    `Expected Content-Type for ${url} to be '${expected}', but got ${actual ?? "'unknown'"}`,
  readBufferNotSupported: (url: string) =>
    `Response from ${url} does not support readBodyBuffer.`,
} as const;
