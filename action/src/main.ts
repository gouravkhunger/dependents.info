import * as core from "@actions/core";

import { ERRORS } from "@/errors";
import { processRepo } from "@/parser";
import { validateRepoName } from "@/utils";

export async function run(): Promise<void> {
  try {
    const name = core.getInput("repo");
    if (!validateRepoName(name)) {
      throw new Error(ERRORS.INVALID_REPO_FORMAT);
    }

    core.info(`Extracting dependents of repository ${name}`);

    const dependents = await processRepo(name);
    const output = dependents.flat();

    core.info(
      `Processed ${output.length} public dependents for repository ${name}`,
    );
    core.setOutput("dependents", output);
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
