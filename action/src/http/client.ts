import { type IncomingHttpHeaders } from "node:http";

import { HttpClient } from "@actions/http-client";

import { ERROR } from "@/constants";

export const get = async (url: string): Promise<string> => {
  const client = new HttpClient();
  try {
    const response = await client.get(url);
    if (response.message.statusCode !== 200) {
      throw new Error(ERROR.failedToFetch(url, response.message.statusCode));
    }
    return await response.readBody();
  } finally {
    client.dispose();
  }
};

export const getImageBuffer = async (
  url: string,
): Promise<[Buffer, IncomingHttpHeaders]> => {
  const client = new HttpClient();
  try {
    const response = await client.get(url);
    if (response.message.statusCode !== 200) {
      throw new Error(ERROR.failedToFetch(url, response.message.statusCode));
    }
    const contentType = response.message.headers["content-type"];
    if (!contentType || !contentType.startsWith("image/")) {
      throw new Error(ERROR.contentTypeMismatch(url, "image", contentType));
    }
    if (typeof response.readBodyBuffer === "undefined") {
      throw new Error(ERROR.readBufferNotSupported(url));
    }
    return [await response.readBodyBuffer(), response.message.headers];
  } finally {
    client.dispose();
  }
};
