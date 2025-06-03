import { processRepo } from "@/parser";
import { parseDependentsPage, parseTotalDependents } from "@/parser/parse";

describe("processRepo function", () => {
  it("should process a repository's dependents correctly", async () => {
    const repo = "gouravkhunger/jekyll-auto-authors";
    const data = await processRepo(repo);
    expect(data.total).toBeGreaterThan(0);
  });
});

describe("parseDependentsPage function", () => {
  it("should parse empty dependents from empty HTML", () => {
    const html = "";
    const { dependents, nextPageLink } = parseDependentsPage(html);
    expect(dependents).toHaveLength(0);
    expect(nextPageLink).toBeUndefined();
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
    expect(dependents).toHaveLength(0);
    expect(nextPageLink).toBe(
      "https://github.com/owner/repo/network/dependents?dependents_after=someString",
    );
  });

  it("should parse next page link when both previous and next buttons are present", () => {
    const html = `
      <div class="paginate-container">
        <div class="BtnGroup" data-test-selector="pagination">
        <a rel="nofollow" class="btn BtnGroup-item" href="https://github.com/owner/repo/network/dependents?dependents_before=next">Previous</a>
        <a rel="nofollow" class="btn BtnGroup-item" href="https://github.com/owner/repo/network/dependents?dependents_after=prev">Next</a></div>
      </div>
    `;
    const { nextPageLink } = parseDependentsPage(html);
    expect(nextPageLink).toBe(
      "https://github.com/owner/repo/network/dependents?dependents_after=prev",
    );
  });

  it("should set nextPageLink to undefined when no next page is available", () => {
    const html = `
      <div class="paginate-container">
        <div class="BtnGroup" data-test-selector="pagination">
          <a rel="nofollow" class="btn BtnGroup-item" href="https://github.com/owner/repo/network/dependents?dependents_before=next">Previous</a>
          <button class="btn BtnGroup-item" disabled="disabled">Next</button></div>
      </div>
    `;
    const { nextPageLink } = parseDependentsPage(html);
    expect(nextPageLink).toBeUndefined();
  });

  it("should set nextPageLink to undefined when no links are present", () => {
    const html = `
      <div class="paginate-container">
        <div class="BtnGroup" data-test-selector="pagination">
          <button class="btn BtnGroup-item" disabled="disabled">Previous</button>
          <button class="btn BtnGroup-item" disabled="disabled">Next</button></div>
      </div>
    `;
    const { nextPageLink } = parseDependentsPage(html);
    expect(nextPageLink).toBeUndefined();
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
        repo: "repo",
        owner: "owner",
        image: "https://avatars.githubusercontent.com/u/46792249?v=4",
      },
    ]);
    expect(nextPageLink).toBe(undefined);
  });
});

describe("parseTotalDependents function", () => {
  it("should parse total dependents from GitHub dependents page HTML snippet", async () => {
    const html = `
      <div role="status" class="table-list-header-toggle states flex-auto pl-0">
        <a class="btn-link selected" href="/owner/repo/network/dependents?dependent_type=REPOSITORY">
          <svg aria-hidden="true" height="16" viewbox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-code-square">
            <path d="M0 1.75C0 .784.784 0 1.75 0h12.5C15.216 0 16 .784 16 1.75v12.5A1.75 1.75 0 0 1 14.25 16H1.75A1.75 1.75 0 0 1 0 14.25Zm1.75-.25a.25.25 0 0 0-.25.25v12.5c0 .138.112.25.25.25h12.5a.25.25 0 0 0 .25-.25V1.75a.25.25 0 0 0-.25-.25Zm7.47 3.97a.75.75 0 0 1 1.06 0l2 2a.75.75 0 0 1 0 1.06l-2 2a.749.749 0 0 1-1.275-.326.749.749 0 0 1 .215-.734L10.69 8 9.22 6.53a.75.75 0 0 1 0-1.06ZM6.78 6.53 5.31 8l1.47 1.47a.749.749 0 0 1-.326 1.275.749.749 0 0 1-.734-.215l-2-2a.75.75 0 0 1 0-1.06l2-2a.751.751 0 0 1 1.042.018.751.751 0 0 1 .018 1.042Z"></path>
          </svg>
          1,364
                          Repositories
        </a>
        <a class="btn-link " href="/owner/repo/network/dependents?dependent_type=PACKAGE">
          <svg aria-hidden="true" height="16" viewbox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-package">
            <path d="m8.878.392 5.25 3.045c.54.314.872.89.872 1.514v6.098a1.75 1.75 0 0 1-.872 1.514l-5.25 3.045a1.75 1.75 0 0 1-1.756 0l-5.25-3.045A1.75 1.75 0 0 1 1 11.049V4.951c0-.624.332-1.201.872-1.514L7.122.392a1.75 1.75 0 0 1 1.756 0ZM7.875 1.69l-4.63 2.685L8 7.133l4.755-2.758-4.63-2.685a.248.248 0 0 0-.25 0ZM2.5 5.677v5.372c0 .09.047.171.125.216l4.625 2.683V8.432Zm6.25 8.271 4.625-2.683a.25.25 0 0 0 .125-.216V5.677L8.75 8.432Z"></path>
          </svg>
          709
                          Packages
        </a>
        <details class="details-reset d-inline-block details-overlay js-dropdown-details position-relative">
          <summary aria-label="Warning" class="d-block px-1">
            <svg aria-hidden="true" height="16" viewbox="0 0 16 16" version="1.1" width="16" data-view-component="true" class="octicon octicon-info">
              <path d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8Zm8-6.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13ZM6.5 7.75A.75.75 0 0 1 7.25 7h1a.75.75 0 0 1 .75.75v2.75h.25a.75.75 0 0 1 0 1.5h-2a.75.75 0 0 1 0-1.5h.25v-2h-.25a.75.75 0 0 1-.75-.75ZM8 6a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z"></path>
            </svg>
          </summary>
          <div class="Popover mt-2 right-0 mr-n2">
            <div class="Popover-message Popover-message--large Box color-shadow-large p-3 Popover-message--top-right ws-normal">
              These counts are approximate and may not exactly match the dependents shown below.
            </div>
          </div>
        </details>
      </div>
    `;
    const totalDependents = parseTotalDependents(html, "owner/repo");
    expect(totalDependents).toBe(1364);
  });
});
