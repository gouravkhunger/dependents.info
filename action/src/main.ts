import { writeFile } from "node:fs/promises";
import path from "node:path";

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

    core.info(
      MESSAGE.processedDependents("public", data.dependents.length, name),
    );
    core.info(MESSAGE.processedDependents("total", data.total, name));
    core.setOutput("dependents", data);

    const distFile = path.join(__dirname, "..", "dist", "dependents.json");
    writeFile(distFile, JSON.stringify(data, null, 2))
      .then(() => core.info(MESSAGE.wroteFile(distFile)))
      .catch((error) => {
        core.error(ERROR.failedToWriteFile(distFile, error.message));
      });
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
