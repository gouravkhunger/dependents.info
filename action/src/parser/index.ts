import * as core from "@actions/core";

import { MESSAGE } from "@/constants";
import { get } from "@/http/client";
import { parseDependentsPage, parseTotalDependents } from "@/parser/parse";
import type { Dependents, ProcessedDependents } from "@/types";
import { buildDependentsUrl } from "@/utils";

export async function processRepo(name: string): Promise<ProcessedDependents> {
  const dependents: Dependents = [];
  let total: number | undefined = undefined;
  let pageLink: string | undefined = buildDependentsUrl(name);

  let count = 0;
  while (pageLink) {
    count++;
    const response = await get(pageLink);
    const page = parseDependentsPage(response);
    dependents.push(...page.dependents);
    pageLink = page.nextPageLink;
    core.info(MESSAGE.processedPage(count, name));
    if (typeof total === "undefined") {
      total = parseTotalDependents(response, name);
    }
  }

  return {
    total: total ?? dependents.length,
    dependents,
  };
}
