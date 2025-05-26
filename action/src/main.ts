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

    const dependents = await processRepo(name);
    const output = dependents.flat();

    core.info(MESSAGE.processedDependents(output.length, name));
    core.setOutput("dependents", output);
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
