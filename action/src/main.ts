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
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
