---
created: 2020-06-17
---

# The Fuck

Наверняка, многие слышали (но не лишним будет напомнить) про утилиту thefuck – «великолепный» инструмент для исправления ошибок в предыдущей команде (в терминале).

https://github.com/nvbn/thefuck (Python, MIT)

<video width="732" height="410" controls>
  <source src="thefuck.mp4" type="video/mp4">
</video>

Мой любимый пример:

```bash
git push
fatal: The current branch master has no upstream branch.
To push the current branch and set the remote as upstream, use

    git push --set-upstream origin master


$ fuck --yeah
git push --set-upstream origin branch_name
Total 0 (delta 0), reused 0 (delta 0)
...
```

Довольно большой список какие команды thefuck умеет исправлять: https://github.com/nvbn/thefuck#how-it-works

#cli