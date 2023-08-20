---
date: 2023-08-19
refs:
  - blog/2020/project-org_ru.md
  - blog/2020/gclone_ru.md
---

# gonew

В Go недавно появился новый экспериментальный инструмент для создания новых проектов: `gonew`:

https://go.dev/blog/gonew

[Шаблоны от Google](https://github.com/GoogleCloudPlatform/go-templates): `appengine`, `pubsubfn`, `httpfn` and `microservice`

[Шаблоны от меня](https://github.com/chuhlomin/gonew): `sever`, `library` и `action`.

Теперь новые проекты я зачинаю с помощью алиаса [`new`](https://github.com/chuhlomin/aliases/blob/main/new.sh):

```sh
cd ~/Projects/chulomin
new server test
```

И новый проект создастся в директории `~/Projects/chulomin/test` с модулем Go, названным `github.com/chuhlomin/test` и всем необходимым старта для серверного проекта.

#go #project
