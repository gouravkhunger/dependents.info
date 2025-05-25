import * as core from "@actions/core";

import { ERRORS } from "@/errors";
import { get } from "@/http/client";
import { parseDependentsPage } from "@/parser/parse";
import { type Dependents } from "@/types";
import { dependentsUrl, validateRepoName } from "@/utils";

export async function run(): Promise<void> {
  try {
    const name = core.getInput("repo");
    if (!validateRepoName(name)) {
      throw new Error(ERRORS.INVALID_REPO_FORMAT);
    }

    core.info(`Extracting dependents of repository ${name}`);
    const dependents: Dependents[] = [];
    let pageLink: string | undefined = dependentsUrl(name);

    while (pageLink) {
      const response = await get(pageLink);
      const page = parseDependentsPage(response);
      dependents.push(page.dependents);
      pageLink = page.nextPageLink;
    }

    core.setOutput("dependents", dependents.flat());
  } catch (error) {
    if (error instanceof Error) {
      core.setFailed(error.message);
    }
  }
}
