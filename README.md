# dependents.info

[![GitHub Network Dependents Count](https://dependents.info/gouravkhunger/dependents.info/badge)](https://dependents.info/gouravkhunger/dependents.info)

easily generate an image of github network dependents to showcase in your project's `readme.md` file.

simply add a github action to your repository and use the image link for your repo.

## demo

here's a demo of the generated image for the gem [`jekyll-auto-authors`](https://github.com/gouravkhunger/jekyll-auto-authors):

<a href="https://dependents.info/gouravkhunger/jekyll-auto-authors">
  <img src="https://dependents.info/gouravkhunger/jekyll-auto-authors/image" />
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
    upload-artifacts: true
    package-id: UGFja2SomeStringzcyMDE
```

| option             | type      | description                                                               | default |
|--------------------|-----------|---------------------------------------------------------------------------|---------|
| `max-pages`        | `number`  | number of network dependents pages to process (max: 100).                 | `50`    |
| `package-id`       | `string`  | use if repo hosts [multiple packages](#multiple-packages). action processes only one at a time. | `""`    |
| `unique-owners`    | `boolean` | disables duplicate users in the generated image.                          | `true`  |
| `exclude-owner`    | `boolean` | exclude repos from the same owner that depend on this repository.         | `true`  |
| `upload-artifacts` | `boolean` | whether to upload the outputs as action's build artifacts.                | `true`  |

#### why github action?

the github action does the heavy lifting of fetching the dependents from your repository's network dependents page.

doing it in a github action makes it much easier to do that from their hosted runners, avoid ip bans, and adhere to the purpose of "archival" of public information as per the [tos](https://docs.github.com/en/site-policy/acceptable-use-policies/github-acceptable-use-policies#7-information-usage-restrictions).

the permission `id-token` with value `write` is required for the action to request a [github oidc](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect) token at the runtime which is then sent to the backend along with the scraped data.

while the `id-token` permission is required to request the oidc token, it in itself is not a security concern to your repository. here's more info from the [github docs](https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect#adding-permissions-settings):

> You won't be able to request the OIDC JWT ID token if the permissions for `id-token` is not set to `write`, however this value doesn't imply granting write access to any resources, only being able to fetch and set the OIDC token for an action or step to enable authenticating with a short-lived access token.

the backend uses the token to verify the `repository` claim directly from github to compare where the data is coming from. mismatched fields will fail the request.

this check ensures that the data backend accepts from a repository comes from it's github action itself. only the action can alter the data in production.

### embed image

> **note**: the image is only available for repositories that run the action successfully.

copy the following code snippet and **replace `owner/repo`** with your repository's name. paste it wherever you want to embed the image.

```html
<a href="https://dependents.info/owner/repo">
  <img src="https://dependents.info/owner/repo/image" />
</a>
```

if you've used the `package-id` option in the action, this should be:

```html
<a href="https://dependents.info/owner/repo?id=idHere">
  <img src="https://dependents.info/owner/repo/image?id=idHere" />
</a>
```

### embed badge

> **note**: the badge is only available for repositories that run the action successfully.

copy the following code snippet and **replace `owner/repo`** with your repository's name. paste it wherever you want to embed the badge.

```html
<a href="https://dependents.info/owner/repo">
  <img src="https://dependents.info/owner/repo/badge" />
</a>
```

if you've used the `package-id` option in the action, this should be:

```html
<a href="https://dependents.info/owner/repo?id=idHere">
  <img src="https://dependents.info/owner/repo/image?id=idHere" />
</a>
```

available query params (optional):

- `logo`: icon name from [simple-icons](https://simpleicons.org).
- `label`: override the default label "dependents".
- `color`: hex, rgb, rgba, hsl, hsla or css named color.
- `logoColor`: hex, rgb, rgba, hsl, hsla or css named color.
- `labelColor`: hex, rgb, rgba, hsl, hsla or css named color.
- `style`: [`flat` (default), `flat-square`, `plastic`, `for-the-badge`, `social`]

usage: `/badge?color=red&style=flat-square`

the badge and the image are self updating so when the github action submits new data, it will be reflected in the readme automatically.

> **note**: in addition to cloudflare's cache lasting up to a day, the image could be cached by github for an extended 7 day period. please refer to [the docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/about-anonymized-urls#removing-an-image-from-camos-cache) on how to manually purge them if required.

## multiple packages

a github repository can host multiple packages, which should be defined with the `package-id` option in the action. other packages can be processed by adding new steps in the same workflow file.

the `package-id` string can be found by going to repository's home page > insights > dependency graph > dependents > select package > copy the text after `?package_id=` in the url.

the same package id should be added to the end of every embedded dependents.info url in the readme with `?id=idHere`.

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

the [makefile](https://github.com/gouravkhunger/dependents.info/blob/main/Makefile) is a convenient way to run the project locally.

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
