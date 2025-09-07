---
date: 2025-09-07
refs:
  - blog/2023/conventional-commits.md
---

# Conventional Comments

![Contentional Comments](conventional-comments.png)

Conventioal Comments is a lightweight convention for writing comments of any kind
(code review comments, documents, etc.) in a consistent and constructive manner.

https://conventionalcomments.org/

Format is as follows:

```
<label> [decorations]: <subject>

[discussion]
```

where decorations and discussion are optional.

Exmaples of comments that follow this convention:

- suggestion(non-blocking): rename variable to `user_id`

  This will make it consistent with other parts of the codebase.

- question: Why are we using a linked list here?
- issue: This function has a time complexity of O(n^2).
- nitpick: Consider using `const` instead of `let` for this variable.
- praise: Great job on implementing the new feature!
- todo: We need to add error handling here.
- note: This part of the code is critical for performance.
- typo: "recieve" should be "receive".
