import * as core from "@actions/core";

import { MESSAGE } from "@/constants";
import { get } from "@/http/client";
import { parseDependentsPage } from "@/parser/parse";
import { type Dependents } from "@/types";
import { dependentsUrl } from "@/utils";

export async function processRepo(name: string): Promise<Dependents[]> {
  const dependents: Dependents[] = [];
  let pageLink: string | undefined = dependentsUrl(name);

  let count = 0;
  while (pageLink) {
    count++;
    const response = await get(pageLink);
    const page = parseDependentsPage(response);
    dependents.push(page.dependents);
    pageLink = page.nextPageLink;
    core.info(MESSAGE.processedPage(count, name));
    if (count >= 5) {
      core.warning(MESSAGE.REACHED_MAX_PAGES);
      break;
    }
  }

  return dependents;
}
