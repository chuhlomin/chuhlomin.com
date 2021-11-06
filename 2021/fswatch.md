---
date: 2021-10-22
---

# fswatch

Sometimes you need to regularly run a program in the terminal when you change files in some directory. This is a kind of "local CI" that frontend engineers are familiar with: change JavaScript file, and the page will update itself in browser ("hot reload").

The program `fswatch` written by Enrico Maria Crisostomo prints the paths of changed files in a given directory, which you can "pipe" to another program.

https://emcrisostomo.github.io/fswatch/

In the example below, `fswatch` runs `genblog` if any file in the current directory has been changed (except for files in directories `.git` and `output`):

```bash
fswatch --one-per-batch --recursive --exclude="output" --exclude=".git" . | xargs -n1 sh -c "genblog"
```

#cli
