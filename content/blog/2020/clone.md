---
date: 2020-07-09
refs:
  - blog/2020/project-org.md
---

# clone

Presenting you the `clone` script, which clones the repository to the folder
following convention `<owner>/<repo>` (like go get before introducing modules):

```bash
clone git@github.com:chuhlomin/terraform.git
↓
mkdir ~/Projects/chuhlomin
git clone git@github.com:chuhlomin/terraform.git ~/Projects/chuhlomin
```

https://github.com/chuhlomin/aliases/blob/main/clone.sh

One more thing to add is the `group_by_owner` script, which arranges existing repositories in the "right" places:

https://gist.github.com/chuhlomin/93770d500e8d25c3fe604d80f3a9abe3

#git
