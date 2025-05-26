import { getImageBuffer } from "@/http/client";

export const buildDependentsUrl = (repoName: string): string => {
  if (!validateRepoName(repoName)) {
    throw new Error("Invalid repository name format");
  }
  return `https://github.com/${repoName}/network/dependents`;
};

export const validateRepoName = (repoName: string): boolean => {
  const regex = /^[a-zA-Z0-9-]+\/[a-zA-Z0-9._-]+$/;
  return regex.test(repoName);
};

export const removeQueryParams = (url: string, ...params: string[]): string => {
  const urlObj = new URL(url);
  params.forEach((param) => urlObj.searchParams.delete(param));
  return urlObj.toString();
};

export const imageUrlToBase64 = async (url: string): Promise<string> => {
  const [buffer, headers] = await getImageBuffer(url);
  const base64 = Buffer.from(buffer).toString("base64");
  const contentType = headers["content-type"] || "image/png";
  return `data:${contentType};base64,${base64}`;
};
