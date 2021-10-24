---
date: 2020-07-09
---

# gclone

Принес вам скрипт gclone, который всегда клонирует репозиторий в нужную папку (как go get до введения модулей):

```bash
gclone git@github.com:chuhlomin/terraform.git
↓
mkdir ~/Projects/chuhlomin
git clone git@github.com:chuhlomin/terraform.git ~/Projects/chuhlomin
```

https://gist.github.com/chuhlomin/f7648f173046251db209585ea68a21ac

Ещё вдогонку скрипт group_by_owner, который раскладывает существующие репозитории по «правильным» местам:

https://gist.github.com/chuhlomin/93770d500e8d25c3fe604d80f3a9abe3

#git
