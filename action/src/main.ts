import { writeFile } from "node:fs/promises";
import path from "node:path";

import { DefaultArtifactClient } from "@actions/artifact";
import * as core from "@actions/core";

import { API_BASE_URL, ERROR, MESSAGE } from "@/constants";
import { processRepo } from "@/parser";
import { validateRepoName } from "@/utils";

export async function run(): Promise<void> {
  try {
    const name = process.env.GITHUB_REPOSITORY ?? "";

    if (!validateRepoName(name)) {
      throw new Error(ERROR.INVALID_REPO_FORMAT);
    }

    core.info(MESSAGE.initExtraction(name));

    const data = await processRepo(name);
    const json = JSON.stringify(data, null, 2);

    core.info(MESSAGE.dependentsCount(data.total, name));
    core.setOutput("dependents", data);

    const distDir = path.join(__dirname, "..", "dist");
    const distFile = path.join(distDir, "dependents.json");

    await writeFile(distFile, json)
      .then(() => core.info(MESSAGE.wroteFile(distFile)))
      .catch((error) => {
        core.error(ERROR.failedToWriteFile(distFile, error.message));
      });

    const uploadArtifacts = core.getInput("upload-artifacts") === "true";
    if (uploadArtifacts) {
      core.info(MESSAGE.artifactUploadLog("started", "dependents.json"));
      const artifact = new DefaultArtifactClient();
      await artifact
        .uploadArtifact("dependents.json", [distFile], distDir)
        .then(() =>
          core.info(MESSAGE.artifactUploadLog("succeeded", "dependents.json")),
        )
        .catch((error) => {
          core.error(ERROR.failedToWriteFile(distFile, error.message));
        });
    }

    let token: string | undefined;
    if (process.env.GITHUB_ACTIONS === "true") {
      token = await core.getIDToken(API_BASE_URL);
    }
    const resp = await fetch(`${API_BASE_URL}/${name}/ingest`, {
      headers: {
        "Content-Type": "application/json",
        ...((token && { Authorization: `Bearer ${token}` }) || {}),
      },
      body: json,
      method: "POST",
    });

    if (!resp.ok) {
      const contentType = resp.headers.get("Content-Type");
      if (contentType && contentType.includes("application/json")) {
        const error = (await resp.json()) as { message: string };
        throw new Error(ERROR.failedToSubmitData(resp.status, error.message));
      } else {
        const text = await resp.text();
        core.error(ERROR.failedToSubmitData(resp.status, text));
        throw new Error(ERROR.failedToSubmitData(resp.status, "unknown error"));
      }
    }

    core.info(MESSAGE.DONE);
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
