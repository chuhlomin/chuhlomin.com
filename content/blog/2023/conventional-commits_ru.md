---
date: 2023-03-11
image: conventional-commits_og.png
refs:
  - blog/2020/gitmoji_ru.md
---

# Соглашение о коммитах

![Conventional Commits](conventional-commits.png)

Спецификация «Соглашение о коммитах» — простое соглашение о том, как нужно писать сообщения коммитов. Оно описывает простой набор правил для создания понятной истории коммитов, а также позволяет проще разрабатывать инструменты автоматизации, основанные на истории коммитов.

https://www.conventionalcommits.org/ru/v1.0.0/

Сообщение коммита должно содержать префикс, который указывает на тип изменений, например:

- feat: add support for Typescript
- fix(gha): release workflow missing secret
- docs: add docs for account service
- chore: update dependencies
- refactor: extract method convertImage

Могу порекомендовать [action-semantic-pull-request](https://github.com/amannn/action-semantic-pull-request)
GitHub Action чтобы автоматически проверять соответствие названия PR этому соглашению.

Это можно настроить для всех PR в GitHub организации с помощью
[Repository Rulesets](https://docs.github.com/en/enterprise-cloud@latest/repositories/configuring-branches-and-merges-in-your-repository/managing-rulesets/about-rulesets).

#git #github
