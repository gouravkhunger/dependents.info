import "./style.css"

import yaml from "@shikijs/langs/yaml";
import html from "@shikijs/langs/html";
import githubDark from "@shikijs/themes/github-dark";

import { createHighlighterCore } from "shiki/core";
import { createOnigurumaEngine } from "shiki/engine/oniguruma";

import { addCopyButton } from "shiki-transformer-copy-button";

const badgeCodeElement = document.querySelector("#badge-code")!;
const embedCodeElement = document.querySelector("#embed-code")!;
const repoInputElement = document.querySelector<HTMLInputElement>("#repo-input")!;
const packageIdInputElement = document.querySelector<HTMLInputElement>("#package-id")!;

const githubAction = `name: Dependents Action

on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  dependents:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
    steps:
      - uses: gouravkhunger/dependents.info@main
`
const actionConfiguration = `- uses: gouravkhunger/dependents.info@main
  with:
    max-pages: 50
    force-run: false
    unique-owners: true
    exclude-owner: true
    upload-artifacts: true
    exclude-users: "user1, user2"
    package-id: UGFja2SomeStringzcyMDE`;

const htmlContent = (name?: string, id?: string) => `<a href="https://dependents.info/${name || "owner/repo"}${id ? `?id=${id}` : ""}">
  <img src="https://dependents.info/${name || "owner/repo"}/image${id ? `?id=${id}` : ""}" />
</a>

Made with [dependents.info](https://dependents.info).`;

const mdContent = (name?: string, id?: string) =>`<a href="https://dependents.info/${name || "owner/repo"}${id ? `?id=${id}` : ""}">
  <img src="https://dependents.info/${name || "owner/repo"}/badge${id ? `?id=${id}` : ""}" />
</a>

Made with [dependents.info](https://dependents.info).`;

const highlighter = await createHighlighterCore({
  langs: [yaml, html],
  themes: [githubDark],
  engine: createOnigurumaEngine(import("shiki/wasm"))
});

const actionCode = highlighter.codeToHtml(githubAction, {
  lang: "yaml",
  theme: "github-dark",
  transformers: [addCopyButton()],
});
document.querySelector("#gh-action")!.innerHTML = actionCode;

const configurationCode = highlighter.codeToHtml(actionConfiguration, {
  lang: "yaml",
  theme: "github-dark",
  transformers: [addCopyButton()],
});
document.querySelector("#action-configuration")!.innerHTML = configurationCode;

const embedCode = (name?: string, id?: string) => highlighter.codeToHtml(htmlContent(name, id), {
  lang: "html",
  theme: "github-dark",
  transformers: [addCopyButton()],
});
embedCodeElement.innerHTML = embedCode();

const badgeCode = (name?: string, id?: string) => highlighter.codeToHtml(mdContent(name, id), {
  lang: "html",
  theme: "github-dark",
  transformers: [addCopyButton()],
});
badgeCodeElement.innerHTML = badgeCode();

const isInvalid = (repo: string): boolean => {
  return !/^[a-zA-Z0-9-]+\/[a-zA-Z0-9._-]+$/.test(repo)
}

repoInputElement.addEventListener("input", (e) => {
  const input = e.target as HTMLInputElement;
  const id = packageIdInputElement.value;
  const name = input.value.trim();
  input.dataset.state = "";
  if (isInvalid(name)) {
    badgeCodeElement.innerHTML = badgeCode(undefined, id);
    embedCodeElement.innerHTML = embedCode(undefined, id);
    if (name === "") {
      input.dataset.state = "";
    } else {
      input.dataset.state = "error";
    }
    return;
  } else {
    input.dataset.state = "";
  }
  input.dataset.state = "valid";
  badgeCodeElement.innerHTML = badgeCode(name, id);
  embedCodeElement.innerHTML = embedCode(name, id);
});

packageIdInputElement.addEventListener("input", (e) => {
  const input = e.target as HTMLInputElement;
  const repo = repoInputElement.value;
  const invalid = isInvalid(repo);
  const id = input.value.trim();

  badgeCodeElement.innerHTML = badgeCode(invalid ? undefined : repo, id);
  embedCodeElement.innerHTML = embedCode(invalid ? undefined : repo, id);
});

const usedByRepos = [
  "joshnuss/svelte-persisted-store",
  "sieblyio/binance",
  "oekazuma/svelte-meta-tags",
  "okineadev/vitepress-plugin-llms",
  "SSShooter/mind-elixir-core",
  "HatsuneMiku3939/direnv-action",
];

const usedBy = document.querySelector("#used-by");
const usedBySection = document.querySelector<HTMLElement>("#used-by-section");

if (usedBy && usedBySection) {
  const rows = await Promise.all(
    usedByRepos.map(async (repo) => {
      try {
        const res = await fetch(`https://dependents.info/${repo}.json`);
        if (!res.ok) {
          return null;
        }
        const data = (await res.json()) as { owner?: string; repo?: string; total?: number };
        const owner = data.owner ?? repo.split("/")[0];
        const name = data.repo ?? repo.split("/")[1];
        const total = typeof data.total === "number" ? data.total : 0;
        return { owner, name, total };
      } catch {
        return null;
      }
    }),
  );

  for (const row of rows) {
    if (!row) {
      continue;
    }
    const slug = `${row.owner}/${row.name}`;
    const a = document.createElement("a");
    a.className = "used-by-card";
    a.href = `/${slug}`;

    const img = document.createElement("img");
    img.src = `https://github.com/${row.owner}.png?size=64`;
    img.width = 32;
    img.height = 32;
    img.alt = "";

    const meta = document.createElement("div");
    meta.className = "used-by-meta";

    const title = document.createElement("span");
    title.className = "used-by-name";
    title.textContent = slug;

    const stars = document.createElement("span");
    stars.className = "used-by-stars";
    stars.textContent = `${row.total.toLocaleString("en-US")} dependents`;

    meta.append(title, stars);
    a.append(img, meta);
    usedBy.append(a);
  }

  if (usedBy.childElementCount > 0) {
    usedBySection.hidden = false;
  }
}
