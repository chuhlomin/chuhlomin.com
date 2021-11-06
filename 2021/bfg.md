---
date: 2021-02-21
---

# BFG Repo-Cleaner

I hope everyone knows that storing secret values in an open form in Git repositories (even private ones) is strongly discouraged. There have been a couple of times when I've cleaned up such repositories. Of course, you have to rewrite the history and coordinate it with the team. The official documentation suggests using git-filter-branch, but there is a faster way.

Roberto Tyley wrote a BFG app in Scala in 2012 that is 10-720 times faster than git-filter-branch.

https://rtyley.github.io/bfg-repo-cleaner/

```bash
git clone --mirror git://example.com/some-big-repo.git
java -jar bfg.jar \
    --strip-blobs-bigger-than 100M \
    --replace-text banned.txt
    some-big-repo.git
cd some-big-repo.git
git reflog expire --expire=now --all && \
    git gc --prune-now--aggressive
git push --force
```

The application knows how to delete:

- passwords, any other secret values
- large files.

Before using it, make sure that all this already absent in the current version of the files.

https://github.com/rtyley/bfg-repo-cleaner

#app #cli #git
