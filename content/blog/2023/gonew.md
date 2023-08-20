---
date: 2023-08-19
refs:
  - blog/2020/project-org.md
  - blog/2020/gclone.md
---

# gonew

Go recently got a new experimental tool for creating new projects: `gonew`:

https://go.dev/blog/gonew

[Templates from Google](https://github.com/GoogleCloudPlatform/go-templates): `appengine`, `pubsubfn`, `httpfn` and `microservice`.

[Templates from me](https://github.com/chuhlomin/gonew): `sever`, `library` and `action`.

Now I start new projects with an [alias `new`](https://github.com/chuhlomin/aliases/blob/main/new.sh):

```sh
cd ~/Projects/chulomin
new server test
```

And a new project will be created in the `~/Projects/chulomin/test` directory with a Go module named `github.com/chuhlomin/test` and everything I need to start a server project.

#go #project
