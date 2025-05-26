import { processRepo } from "@/parser";
import { parseDependentsPage } from "@/parser/parse";

describe("processRepo function", () => {
  it("should process a repository's dependents correctly", async () => {
    const repo = "gouravkhunger/jekyll-auto-authors";
    const dependents = await processRepo(repo);
    expect(dependents.length).toBeGreaterThan(0);
  });
});

describe("parseDependentsPage function", () => {
  it("should parse empty dependents from empty HTML", () => {
    const html = "";
    const { dependents, nextPageLink } = parseDependentsPage(html);
    expect(dependents).toEqual([]);
    expect(nextPageLink).toBe("");
  });

  it("should parse next page link from GitHub dependents page HTML snippet", () => {
    const html = `
      <div class="paginate-container">
        <div class="BtnGroup" data-test-selector="pagination">
          <button class="btn BtnGroup-item" disabled="disabled">Previous</button>
          <a rel="nofollow" class="btn BtnGroup-item" href="https://github.com/owner/repo/network/dependents?dependents_after=someString">Next</a>
        </div>
      </div>
    `;
    const { dependents, nextPageLink } = parseDependentsPage(html);
    expect(dependents).toEqual([]);
    expect(nextPageLink).toBe(
      "https://github.com/owner/repo/network/dependents?dependents_after=someString",
    );
  });

  it("should parse dependents from GitHub dependents page HTML", () => {
    const html = `
      <div class="Box-row d-flex flex-items-center" data-test-id="dg-repo-pkg-dependent">
        <img class="avatar mr-2 avatar-user" src="https://avatars.githubusercontent.com/u/46792249?s=40&v=4" width="20" height="20" alt="@owner" />

        <span class="f5 color-fg-muted" data-repository-hovercards-enabled>
          <a data-hovercard-type="user" data-hovercard-url="/users/owner/hovercard" data-octo-click="hovercard-link-click" data-octo-dimensions="link_type:self" href="/owner">owner</a> /
          <a class="text-bold" data-hovercard-type="repository" data-hovercard-url="/owner/repo/hovercard" href="/owner/repo">repo</a>
            <small></small>
        </span>
        <div class="d-flex flex-auto flex-justify-end">
          <span class="color-fg-muted text-bold pl-3">
            <svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-star">
        <path d="M8 .25a.75.75 0 0 1 .673.418l1.882 3.815 4.21.612a.75.75 0 0 1 .416 1.279l-3.046 2.97.719 4.192a.751.751 0 0 1-1.088.791L8 12.347l-3.766 1.98a.75.75 0 0 1-1.088-.79l.72-4.194L.818 6.374a.75.75 0 0 1 .416-1.28l4.21-.611L7.327.668A.75.75 0 0 1 8 .25Zm0 2.445L6.615 5.5a.75.75 0 0 1-.564.41l-3.097.45 2.24 2.184a.75.75 0 0 1 .216.664l-.528 3.084 2.769-1.456a.75.75 0 0 1 .698 0l2.77 1.456-.53-3.084a.75.75 0 0 1 .216-.664l2.24-2.183-3.096-.45a.75.75 0 0 1-.564-.41L8 2.694Z"></path>
    </svg>
            28
          </span>
          <span class="color-fg-muted text-bold pl-3">
            <svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-repo-forked">
        <path d="M5 5.372v.878c0 .414.336.75.75.75h4.5a.75.75 0 0 0 .75-.75v-.878a2.25 2.25 0 1 1 1.5 0v.878a2.25 2.25 0 0 1-2.25 2.25h-1.5v2.128a2.251 2.251 0 1 1-1.5 0V8.5h-1.5A2.25 2.25 0 0 1 3.5 6.25v-.878a2.25 2.25 0 1 1 1.5 0ZM5 3.25a.75.75 0 1 0-1.5 0 .75.75 0 0 0 1.5 0Zm6.75.75a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Zm-3 8.75a.75.75 0 1 0-1.5 0 .75.75 0 0 0 1.5 0Z"></path>
    </svg>
            5
          </span>
        </div>
      </div>
    `;

    const { dependents, nextPageLink } = parseDependentsPage(html);
    expect(dependents.length).toBe(1);
    expect(dependents).toEqual([
      {
        stars: 28,
        name: "owner/repo",
        image: "https://avatars.githubusercontent.com/u/46792249?v=4",
      },
    ]);
    expect(nextPageLink).toBe("");
  });
});
