---
date: 2020-05-06
---

# .gitignore

Today, the primary version control system is Git, which has replaced SVN.
https://www.openhub.net/repositories/compare

Microsoft owns GitHub, the leading site for open-source software development.

When working on a project, you often need to create temporary files (`*.out`,
`*.log`, `*.jar`) or files with secrets/passwords (`.env`, `.drone.sec.yml`),
which should not get into the code repository.

Git creates a `.gitignore` file for this purpose in the project directory,
which describes what to ignore. It's also pretty easy to add a `.gitignore`
for all the repositories in your system, where you can put a `.DS_Store`
or `Thumbs.db` (which the operating system creates automatically) once
and for all, so you don't have to drag them from project to project:

```bash
git config --global core.excludesfile '~/.gitignore'
```

By the way, I recently saw a nice collection of `.gitignore` files
for different needs:
https://github.com/github/gitignore

#git
