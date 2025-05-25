import * as core from "@actions/core";

import { ERRORS } from "@/errors";
import { validateRepoName } from "@/utils";

export async function run(): Promise<void> {
  try {
    const name = core.getInput("repo");
    if (!validateRepoName(name)) {
      throw new Error(ERRORS.INVALID_REPO_FORMAT);
    }
    core.info(`Repository: ${name}!`);
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
