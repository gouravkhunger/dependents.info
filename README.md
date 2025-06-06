# dependents.info

easily generate an image of github network dependents to showcase in your project's `readme.md` file.

simply add a github action to your repository and use the image link for your repo.

## demo

here's a demo of the generated image for the gem [`jekyll-auto-authors`](https://github.com/gouravkhunger/jekyll-auto-authors):

<a href="https://github.com/gouravkhunger/jekyll-auto-authors/network/dependents">
  <img src="https://dependents.info/gouravkhunger/jekyll-auto-authors/image.svg" />
</a>

## quickstart

### github action

add this file to your repository's `.github/workflows` folder.

`dependents.yml`

```yml
name: Dependents Action

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
```

once you push this file, the action will process the dependents for the repository and the api will generate the image.

> tip: instead of running the action on `push`, you can use a [cron job](https://docs.github.com/en/actions/writing-workflows/choosing-when-your-workflow-runs/events-that-trigger-workflows#schedule) to schedule the action.

#### configuration (optional)

add the following options to your `dependents.yml` file if you want to customize the action's behavior:

```yml
  - uses: gouravkhunger/dependents.info@main
    with:
      max-pages: 50
      unique-owners: true
      exclude-owner: true
      upload-artifacts: false
```

| option             | type      | description                                                                 | default |
|--------------------|-----------|-----------------------------------------------------------------------------|---------|
| `max-pages`        | `number`  | maximum number of network dependents pages to process (max: 100).            | `50`    |
| `unique-owners`    | `boolean` | whether to disable unique users in the generated image.                      | `true`  |
| `exclude-owner`    | `boolean` | whether to exclude repos from the same owner that depend on this repository. | `true`  |
| `upload-artifacts` | `boolean` | whether to upload the outputs as action's build artifacts.                   | `false` |

### embed image

> **note**: the image is only available for repositories that run the action successfully.

copy the following code snippet below and **replace `owner/repo`** with your repository's name. paste it wherever you want to embed the image.

```html
<a href="https://github.com/owner/repo/network/dependents">
  <img src="https://dependents.info/owner/repo/image.svg" />
</a>

Made with [dependents.info](https://dependents.info).
```

the image is self updating so when the github action submits new data, it will be reflected in the readme automatically.

> **note**: in addition to cloudflare's cache lasting up to a day, the image could be cached by github for an extended 7 day period. please refer to [the docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-anonymized-urls#removing-an-image-from-camos-cache) on how to manually purge them if required.

## stack

this project is built as a monorepo hosting three packages in specific folders:

- `www`: the frontend built with [vite](https://vite.dev) using the [vanilla-ts](https://vite.dev/guide/#scaffolding-your-first-vite-project) template.

  it's a simple site that acts as the project home page. built and sent to the api to be served on the home page.

- `action`: the github action built in typescript with [@actions/core](https://github.com/actions/toolkit/tree/main/packages/core).

  it is responsible for processing the repository it is run on and building the dependents data required by the api to generate the image.

  uses [github oidc](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect) to authenticate with the api and thus preventing abuse of the service.

- `api`: [go](https://go.dev) based api built with the [fiber](https://gofiber.io) web framework.

  ingests data from the github action and processes the svg image. uses [badger](https://github.com/hypermodeinc/badger) db to efficiently store the generated images.

## contributing

the project uses [node.js](https://nodejs.org) `>=20` and [go](https://go.dev) version `>=1.24`.

the [makefile](https://github.com/gouravkhunger/dependents.info/blob/main/Makefile) is a convinient way to run the project locally.

```bash
git clone https://github.com/gouravkhunger/dependents.info
cd dependents.info
make install
```

to contribute:

- fork the repository.
- create a new branch:
  ```bash
  git checkout -b branch
  ```
- add test cases for your feature / bug fix.
- add your changes and ensure the tests pass:
  ```bash
  make api-test
  make action-test
  ```
- make a commit:
  ```bash
  git commit -m "<message>"
  ```
  please try to follow the [conventional commits](https://www.conventionalcommits.org) format for messages.
- push & create a pull request.

## inspiration

this project is inspired by the awesome work done in:

- [contributors-img](https://github.com/lacolaco/contributors-img)
- [github-dependents-info](https://github.com/nvuillam/github-dependents-info)

## license

this project is [mit licensed](https://github.com/gouravkhunger/dependents.info/blob/main/LICENSE).
