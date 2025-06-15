import { API_BASE_URL, ERROR } from "@/constants";
import { getImageBuffer } from "@/http/client";

export const buildAPIUrl = (name: string, id: string): string => {
  if (!validateRepoName(name)) {
    throw new Error(ERROR.INVALID_REPO_FORMAT);
  }
  return `${API_BASE_URL}/${name}/ingest${id ? `?id=${id}` : ""}`;
};

export const buildDependentsUrl = (name: string, id: string): string => {
  if (!validateRepoName(name)) {
    throw new Error(ERROR.INVALID_REPO_FORMAT);
  }
  return `https://github.com/${name}/network/dependents${id ? `?package_id=${id}` : ""}`;
};

export const validateRepoName = (name: string): boolean => {
  const regex = /^[a-zA-Z0-9-]+\/[a-zA-Z0-9._-]+$/;
  return regex.test(name);
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
