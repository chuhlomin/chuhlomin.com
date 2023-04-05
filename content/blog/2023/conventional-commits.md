---
date: 2023-03-11
image: conventional-commits_og.png
refs:
  - blog/2020/gitmoji.md
---

# Conventional Commits

![Conventional Commits](conventional-commits.png)

Conventional Commits is a specification for adding human and machine readable meaning to commit messages.
This allows automate the release process and generate changelogs.

https://www.conventionalcommits.org/en/v1.0.0/

Each commit message prefix is a type of change, examples:

* feat: add support for Typescript
* fix(gha): release workflow missing secret
* docs: add docs for account service
* chore: update dependencies
* refactor: extract method convertImage

Can recommend the [action-semantic-pull-request](https://github.com/amannn/action-semantic-pull-request)
GitHub Action to enforce this convention for pull requests titles.

This can be enforced for every PR in a GitHub org by using
[Required Workflows](https://docs.github.com/en/actions/using-workflows/required-workflows).

#git #github
