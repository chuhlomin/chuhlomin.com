---
created: 2020-07-09
---

# gclone

Presenting you the `gclone` script, which clones the repository to the folder
following convention `<owner>/<repo>` (like go get before introducing modules):

```bash
gclone git@github.com:chuhlomin/terraform.git
↓
mkdir ~/Projects/chuhlomin
git clone git@github.com:chuhlomin/terraform.git ~/Projects/chuhlomin
```

https://gist.github.com/chuhlomin/f7648f173046251db209585ea68a21ac

One more thing to add is the `group_by_owner` script, which arranges existing repositories in the "right" places:

https://gist.github.com/chuhlomin/93770d500e8d25c3fe604d80f3a9abe3

#git
