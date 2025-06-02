import * as cheerio from "cheerio";

import { ERROR } from "@/constants";
import { Dependents, DependentsPage } from "@/types";
import { removeQueryParams } from "@/utils";

export const parseTotalDependents = (doc: string, repo: string): number => {
  const $ = cheerio.load(doc);
  const anchor = $(
    `a[href="/${repo}/network/dependents?dependent_type=REPOSITORY"]`,
  );
  const dependents = anchor.text().trim().split(" ")[0];
  if (!dependents) throw new Error(ERROR.failedToParseTotalDependents(repo));
  return parseInt(dependents.replace(/,/g, ""), 10);
};

export const parseDependentsPage = (doc: string): DependentsPage => {
  const $ = cheerio.load(doc);
  const dependents: Dependents = [];
  $('[data-test-id="dg-repo-pkg-dependent"]').each((_, el) => {
    const owner = $(el)
      .find(
        '[data-hovercard-type="user"], [data-hovercard-type="organization"]',
      )
      .text()
      .trim();
    const repo = $(el).find('[data-hovercard-type="repository"]').text().trim();

    const starsText = $(el).find(".octicon-star").parent().text().trim();
    const stars = parseInt(starsText.replace(/,/g, ""), 10) || 0;

    const imageSrc = $(el).find("img").attr("src") || "";
    const image = removeQueryParams(imageSrc, "s");

    dependents.push({ owner, repo, stars, image });
  });

  const nextPageLink =
    $('[data-test-selector="pagination"]').find("a").attr("href") || "";

  return { dependents, nextPageLink };
};
