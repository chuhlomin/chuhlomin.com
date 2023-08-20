---
date: 2020-07-09
refs:
  - blog/2020/project-org_ru.md
---

# clone

Принес вам скрипт gclone, который всегда клонирует репозиторий в нужную папку (как go get до введения модулей):

```bash
clone git@github.com:chuhlomin/terraform.git
↓
mkdir ~/Projects/chuhlomin
git clone git@github.com:chuhlomin/terraform.git ~/Projects/chuhlomin
```

https://github.com/chuhlomin/aliases/blob/main/clone.sh

Ещё вдогонку скрипт group_by_owner, который раскладывает существующие репозитории по «правильным» местам:

https://gist.github.com/chuhlomin/93770d500e8d25c3fe604d80f3a9abe3

#git
