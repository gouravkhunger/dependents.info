import "./style.css"

declare global {
  interface Window {
    embedCode: (repo?: string) => string;
  }
}

const embedCodeElement = document.querySelector<HTMLDivElement>("#embed-code")!;
const repoInputElement = document.querySelector<HTMLInputElement>("#repo-input")!;

repoInputElement.addEventListener("input", (e) => {
  const input = e.target as HTMLInputElement;
  const name = input.value.trim();
  input.dataset.state = "";
  if (!/^[a-zA-Z0-9-]+\/[a-zA-Z0-9._-]+$/.test(name)) {
    embedCodeElement.innerHTML = window.embedCode();
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
  embedCodeElement.innerHTML = window.embedCode(name);
});
