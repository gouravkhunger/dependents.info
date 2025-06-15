import "./style.css"

declare global {
  interface Window {
    embedCode: (repo?: string, id?: string) => string;
    badgeCode: (repo?: string, id?: string) => string;
  }
}

const badgeCodeElement = document.querySelector<HTMLDivElement>("#badge-code")!;
const embedCodeElement = document.querySelector<HTMLDivElement>("#embed-code")!;
const repoInputElement = document.querySelector<HTMLInputElement>("#repo-input")!;
const packageIdInputElement = document.querySelector<HTMLInputElement>("#package-id")!;

const isInvalid = (repo: string): boolean => {
  return !/^[a-zA-Z0-9-]+\/[a-zA-Z0-9._-]+$/.test(repo)
}

repoInputElement.addEventListener("input", (e) => {
  const input = e.target as HTMLInputElement;
  const id = packageIdInputElement.value;
  const name = input.value.trim();
  input.dataset.state = "";
  if (isInvalid(name)) {
    badgeCodeElement.innerHTML = window.badgeCode(undefined, id);
    embedCodeElement.innerHTML = window.embedCode(undefined, id);
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
  badgeCodeElement.innerHTML = window.badgeCode(name, id);
  embedCodeElement.innerHTML = window.embedCode(name, id);
});

packageIdInputElement.addEventListener("input", (e) => {
  const input = e.target as HTMLInputElement;
  const repo = repoInputElement.value;
  const invalid = isInvalid(repo);
  const id = input.value.trim();

  badgeCodeElement.innerHTML = window.badgeCode(invalid ? undefined : repo, id);
  embedCodeElement.innerHTML = window.embedCode(invalid ? undefined : repo, id);
});
