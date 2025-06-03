import * as core from "@actions/core";

import { MESSAGE } from "@/constants";
import { get } from "@/http/client";
import { parseDependentsPage, parseTotalDependents } from "@/parser/parse";
import type { Dependents, ProcessedDependents } from "@/types";
import { buildDependentsUrl, imageUrlToBase64 } from "@/utils";

export async function processRepo(name: string): Promise<ProcessedDependents> {
  const dependents: Dependents = [];
  let total: number | undefined = undefined;
  let pageLink: string | undefined = buildDependentsUrl(name);

  const maxPages = parseInt(core.getInput("max-pages"), 10);
  const safeMaxPages = Math.max(0, Math.min(maxPages, 100));

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
    if (count >= safeMaxPages) {
      core.info(MESSAGE.maxPagesReached(safeMaxPages, name));
      break;
    }
  }

  let data;
  const useUniqueOwners = core.getInput("unique-owners") === "true";
  const sortedData = dependents.sort((a, b) => b.stars - a.stars);

  if (useUniqueOwners) {
    data = Array.from(
      new Map(sortedData.map((obj) => [obj.owner, obj])).values(),
    );
  } else {
    data = sortedData;
  }

  const transformedData = data.slice(0, 10).map(async (dep) => ({
    repo: dep.repo,
    owner: dep.owner,
    image: await imageUrlToBase64(dep.image),
  }));

  return {
    total: total ?? dependents.length,
    dependents: await Promise.all(transformedData),
  };
}
