import { writeFile } from "node:fs/promises";
import path from "node:path";

import { DefaultArtifactClient } from "@actions/artifact";
import * as core from "@actions/core";

import { ERROR, MESSAGE } from "@/constants";
import { processRepo } from "@/parser";
import { validateRepoName } from "@/utils";

export async function run(): Promise<void> {
  try {
    const name = core.getInput("repo");
    if (!validateRepoName(name)) {
      throw new Error(ERROR.INVALID_REPO_FORMAT);
    }

    core.info(MESSAGE.initExtraction(name));

    const data = await processRepo(name);

    core.info(MESSAGE.dependentsCount(data.total, name));
    core.setOutput("dependents", data);

    const distDir = path.join(__dirname, "..", "dist");
    const distFile = path.join(distDir, "dependents.json");
    writeFile(distFile, JSON.stringify(data, null, 2))
      .then(() => core.info(MESSAGE.wroteFile(distFile)))
      .catch((error) => {
        core.error(ERROR.failedToWriteFile(distFile, error.message));
      });

    const uploadArtifacts = core.getInput("upload-artifacts") === "true";
    if (!uploadArtifacts) {
      core.info(MESSAGE.DONE);
      return;
    }

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

    core.info(MESSAGE.DONE);
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
