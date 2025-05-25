import { type IncomingHttpHeaders } from "node:http";

import { HttpClient } from "@actions/http-client";

const client = new HttpClient();

export const get = async (url: string): Promise<string> => {
  try {
    const response = await client.get(url);
    if (response.message.statusCode !== 200) {
      throw new Error(`Failed to fetch ${url}: ${response.message.statusCode}`);
    }
    return await response.readBody();
  } finally {
    client.dispose();
  }
};

export const getImageBuffer = async (
  url: string,
): Promise<[Buffer, IncomingHttpHeaders]> => {
  try {
    const response = await client.get(url);
    if (response.message.statusCode !== 200) {
      throw new Error(`Failed to fetch ${url}: ${response.message.statusCode}`);
    }
    const contentType = response.message.headers["content-type"];
    if (!contentType || !contentType.startsWith("image/")) {
      throw new Error(
        `Content-Type is not an image for ${url}: ${contentType}`,
      );
    }
    if (typeof response.readBodyBuffer === "undefined") {
      throw new Error(`Response does not support readBodyBuffer for ${url}`);
    }
    return [await response.readBodyBuffer(), response.message.headers];
  } finally {
    client.dispose();
  }
};
