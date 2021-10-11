---
created: 2020-05-30
---

# Gitmoji

![Gitmoji promo](gitmoji.png "Gitmoji promo")

Somehow I missed such a thing as Gitmoji (since 2016).
It's the emoji "guide" for commit messages.

https://gitmoji.carloscuesta.me

For example, if a commit fixes some bug, it suggests writing "🐛" prefix to the message.  
Or "🔥" if it removes any code or files. This makes it easy to determine the purpose of the commit by just looking at the emoji used.

It's an open-source project, and 44 people have already contributed to the list.  
Especially nice is that the commits themselves to this file already follow these guidelines:
https://github.com/carloscuesta/gitmoji/blame/master/src/data/gitmojis.json

Note that when you click on an emoji on the site, it gets copied to the clipboard.

And if you don't want to always go to the site - there's `gitmoji-cli`, which serves as a "shell" over the git commit (see screenshot).

https://github.com/carloscuesta/gitmoji-cli (MIT)

#git #emoji #interface
